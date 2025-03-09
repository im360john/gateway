package cli

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/xerrors"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/centralmind/gateway/connectors"
	gw_model "github.com/centralmind/gateway/model"
)

var (
	basePrompt = `
!Important rules:
	- The most important: The final output must contain *only valid single JSON* with no additional commentary, explanations, or markdown formatting!
	- The JSON configuration must strictly adhere to the provided JSON schema, including all required fields.
	- Description of API endpoints should also have an example, to help chatbot to use it.
	- All descriptions and summary must not have any sensitive information/data from security point of view including database types, password and etc.
	- All SQL queries must be Pure SQL that will be used in golang SQLx on top of database - {database_type} and be fully parameterized (using named parameters) to prevent SQL injection.
	- Do not generate output schema for endpoints.
	- All SQL queries must be verified that they will not return array of data where expected one item.
	- SQL queries should be optimized for {database_type} and use appropriate indexes.
	- Endpoints that return lists must include pagination parameters (offset and limit).
	- Consistent Endpoint Definitions: Each table defined in the DDL should have corresponding endpoints as specified by the JSON schema, including method, path, description, SQL query, and parameters.
	- Sensitive Data Handling: If any columns contain sensitive or PII data like phone number, SSN, address, credit card etc, they must be flagged appropriately (e.g., using a "pii" flag).
	- Each Parameter in API endpoints may have default value taken from corresponded example rows, only if it's not PII or sensitive data
	- If some entity requires pagination, there should be separate API that calculates total_count, so pagination can be queried
	- For Postgres, use all table names and column names in double quotes, e.g., "table_name" and "column_name". 
	- If a schema is specified in the table name (format: schema.table), use it in your queries appropriately for the database type. For Postgres, this would be "schema"."table_name".
`
)

var (
	//go:embed api_config_schema.json
	apiConfigSchema []byte
)

type TableData struct {
	Columns  []gw_model.ColumnSchema
	Name     string
	Sample   []map[string]any
	RowCount int
}

// PromptColumnSchema is used specifically for generating prompts,
// omitting sensitive fields like PII flag
type PromptColumnSchema struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	PrimaryKey bool   `yaml:"primary_key,omitempty"`
}

func columnToPromptSchema(col gw_model.ColumnSchema) PromptColumnSchema {
	return PromptColumnSchema{
		Name:       col.Name,
		Type:       string(col.Type),
		PrimaryKey: col.PrimaryKey,
	}
}

func init() {
	// Configure logrus for nicer output
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		ForceColors:            true,
		FullTimestamp:          false,
		TimestampFormat:        "",
		DisableTimestamp:       true,
	})
}

func Discover() *cobra.Command {
	var configPath string
	var databaseType string
	var tables string
	var aiProvider string
	var aiEndpoint string
	var aiAPIKey string
	var aiModel string
	var aiMaxTokens int
	var aiTemperature float32
	var aiReasoning bool
	var bedrockRegion string
	var vertexAIRegion string
	var vertexAIProject string
	var output string
	var extraPrompt string
	var promptFile string
	var openaiLogFile string

	//var red string = "\033[31m"
	//var green = "\033[32m"
	var cyan = "\033[36m"
	var yellow = "\033[33m"
	var violet = "\033[35m"
	var reset = "\033[0m" // reset color

	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover generates gateway config",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			startTime := time.Now()

			// Parse comma-separated tables list
			var tablesList []string
			if tables != "" {
				tablesList = strings.Split(tables, ",")
				// Trim spaces from table names
				for i := range tablesList {
					tablesList[i] = strings.TrimSpace(tablesList[i])
				}
			}

			// Configure header
			logrus.Info("\r\n")
			logrus.Info("🚀 API Discovery Process")

			logrus.Info("Step 1: Read configs")
			configRaw, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			connector, err := connectors.New(databaseType, configRaw)
			if err != nil {
				return err
			}
			if err := connector.Ping(context.Background()); err != nil {
				return err
			}
			logrus.Info("✅ Step 1 completed. Done.")
			logrus.Info("\r\n")

			logrus.Info("Step 2: Discover data")
			allTables, err := connector.Discovery(context.Background())
			if err != nil {
				return err
			}

			tableSet := map[string]bool{}
			for _, table := range tablesList {
				tableSet[table] = true
			}
			if len(tablesList) == 0 {
				for _, table := range allTables {
					tableSet[table.Name] = true
				}
			}

			// Show discovered tables
			logrus.Info("Discovered Tables:")
			for _, table := range allTables {
				if tableSet[table.Name] {
					logrus.Infof("  - "+cyan+"%s"+reset+": "+yellow+"%d"+reset+" columns, "+yellow+"%d"+reset+" rows", table.Name, len(table.Columns), table.RowCount)
				}
			}

			// Check if any tables were found after filtering
			var filteredTablesCount int
			for range tableSet {
				filteredTablesCount++
			}
			if filteredTablesCount == 0 {
				return fmt.Errorf("error: no tables found to process. Please verify your database connection and table selection criteria")
			}

			logrus.Info("✅ Step 2 completed. Done.")
			logrus.Info("\r\n")
			// Sample data
			logrus.Info("Step 3: Sample data from tables")
			var tablesToGenerate []TableData
			for _, table := range allTables {
				if !tableSet[table.Name] {
					continue
				}
				sample, err := connector.Sample(context.Background(), table)
				if err != nil {
					return err
				}
				tablesToGenerate = append(tablesToGenerate, TableData{
					Columns:  table.Columns,
					Name:     table.Name,
					Sample:   sample,
					RowCount: table.RowCount,
				})
			}

			// Show sampled data
			logrus.Info("Data Sampling Results:")
			for _, table := range tablesToGenerate {
				logrus.Infof("  - "+cyan+"%s"+reset+": "+yellow+"%d"+reset+" rows sampled", table.Name, len(table.Sample))
			}

			logrus.Info("✅ Step 3 completed. Done.")
			logrus.Info("\r\n")
			// Prepare prompt
			logrus.Info("Step 4: Prepare the prompt for the AI")
			fullPrompt := generatePrompt(databaseType, extraPrompt, tablesToGenerate, getSchemaFromConfig(databaseType, configRaw))
			if err := saveToFile(promptFile, fullPrompt); err != nil {
				logrus.Error("failed to save prompt:", err)
			}
			logrus.Infof("Prompt saved locally to %s", promptFile)

			logrus.Info("✅ Step 4 completed. Done.")
			logrus.Info("\r\n")
			// Call API
			logrus.Info("Step 5: Using AI to design API")
			config, resp, err := callOpenAI(aiAPIKey, fullPrompt, aiEndpoint, aiModel)
			if err != nil {
				logrus.Error("failed to call OpenAI:", err)
				return err
			}

			// Show generated API endpoints
			var apiEndpoints int
			logrus.Info("API Functions Created:")
			for _, table := range config.Database.Tables {
				for _, endpoint := range table.Endpoints {
					logrus.Infof("  - "+cyan+"%s"+reset+" "+violet+"%s"+reset+" - %s", endpoint.HTTPMethod, endpoint.HTTPPath, endpoint.Summary)
					apiEndpoints++
				}
			}

			config.Database.Type = databaseType
			config.Database.Connection = string(configRaw)

			// Save configuration
			configData, err := yaml.Marshal(config)
			if err != nil {
				logrus.Error("yaml failed:", err)
				return err
			}

			if err := saveToFile(output, string(configData)); err != nil {
				logrus.Error("failed:", err)
				return err
			}
			logrus.Info("\r\n")
			logrus.Infof("API schema saved to: "+cyan+"%s"+reset, output)

			logrus.Info("✅ Step 5: API Specification Generation Completed!")
			logrus.Info("\r\n")
			// Show statistics
			duration := time.Since(startTime)
			tokensUsed := resp.Usage.TotalTokens

			logrus.Info("✅ All steps completed. Done.")
			logrus.Info("\r\n")
			logrus.Info("--- Execution Statistics ---")
			logrus.Infof("Total time taken: "+yellow+"%v"+reset, duration.Round(time.Second))
			logrus.Infof("Tokens used: "+yellow+"%d"+reset+" (Estimated cost: "+violet+"$%.4f"+reset+")",
				tokensUsed, (float64(resp.Usage.PromptTokens)/1000000)*1.1+(float64(resp.Usage.CompletionTokens)/1000000)*4.4)
			logrus.Infof("Tables processed: "+yellow+"%d"+reset, len(tablesToGenerate))
			logrus.Infof("API methods created: "+yellow+"%d"+reset, apiEndpoints)

			// Count PII columns from the generated config
			var piiColumnsCount int
			for _, table := range config.Database.Tables {
				for _, column := range table.Columns {
					if column.PII {
						piiColumnsCount++
					}
				}
			}
			logrus.Infof("Total number of columns with PII data: "+yellow+"%d"+reset, piiColumnsCount)

			answerText := strings.TrimSpace(resp.Choices[0].Message.Content)
			if err := saveToFile(openaiLogFile, answerText); err != nil {
				logrus.Error("failed to save OpenAI response:", err)
			}

			var res gw_model.Config
			if err := yaml.Unmarshal([]byte(answerText), &res); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "connection.yaml", "Path to connection yaml file")
	cmd.Flags().StringVar(&tables, "tables", "", "Comma-separated list of tables to include (e.g. 'table1,table2,table3')")
	cmd.Flags().StringVar(&databaseType, "db-type", "postgres", "Type of database")
	/*
		AI provider options:
	*/
	cmd.Flags().StringVar(&aiProvider, "ai-provider", "openai", "AI provider to use")
	cmd.Flags().StringVar(&aiEndpoint, "ai-endpoint", "", "Custom OpenAI-compatible API endpoint URL")
	cmd.Flags().StringVar(&aiAPIKey, "ai-api-key", "ai-api-key", "AI API token")
	cmd.Flags().StringVar(&bedrockRegion, "bedrock-region", "", "Bedrock region")
	cmd.Flags().StringVar(&vertexAIRegion, "vertexai-region", "", "Vertex AI region")
	cmd.Flags().StringVar(&vertexAIProject, "vertexai-project", "", "Vertex AI project")
	cmd.Flags().StringVar(&aiModel, "ai-model", "", "AI model to use")
	cmd.Flags().IntVar(&aiMaxTokens, "ai-max-tokens", 0, "Maximum tokens to use")
	cmd.Flags().Float32Var(&aiTemperature, "ai-temperature", -1.0, "AI temperature")
	cmd.Flags().BoolVar(&aiReasoning, "ai-reasoning", true, "Enable reasoning")

	cmd.Flags().StringVar(&output, "output", "gateway.yaml", "Resulted yaml path")
	cmd.Flags().StringVar(&extraPrompt, "prompt", "generate reasonable set of API-s for this data", "Custom input to generate API-s")
	cmd.Flags().StringVar(&promptFile, "prompt-file", filepath.Join(getDefaultLogDir(), "prompt_default.txt"), "Path to save the generated prompt")
	cmd.Flags().StringVar(&openaiLogFile, "openai-log", filepath.Join(getDefaultLogDir(), "open-ai-raw.log"), "Path to save OpenAI raw response")

	return cmd
}

func generatePrompt(databaseType, extraPrompt string, tables []TableData, schema string) string {
	res := "I need a config for an automatic API that will be used by another AI bot or LLMs..."
	res += "\n"
	res += strings.ReplaceAll(basePrompt, "{database_type}", databaseType)
	res += "\n" + string(apiConfigSchema) + "\n" + extraPrompt + "\n\n"
	for _, table := range tables {
		// Apply schema to table name if schema is provided and not empty
		var tableName string

		if schema != "" {
			// Qualify the table name with schema
			tableName = fmt.Sprintf("%s.%s", schema, table.Name)
		} else {
			// Use the table name as is
			tableName = fmt.Sprintf("%s.%s", "public", table.Name)
		}

		res += fmt.Sprintf(`
<%[1]s number_columns=%[5]v number_rows=%[6]v>
schema:
%[2]s
---
data_sample:
%[3]s
</%[1]s>

`, tableName, yamlify(table.Columns), yamlify(table.Sample), len(table.Sample), len(table.Columns), table.RowCount)
	}

	return res
}

func yamlify(sample any) string {
	// Convert ColumnSchema to PromptColumnSchema if needed
	if columns, ok := sample.([]gw_model.ColumnSchema); ok {
		promptColumns := make([]PromptColumnSchema, len(columns))
		for i, col := range columns {
			promptColumns[i] = columnToPromptSchema(col)
		}
		sample = promptColumns
	}
	raw, _ := yaml.Marshal(sample)
	return string(raw)
}

func saveToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

// startSpinner starts a loading animation in the console
func startSpinner(message string, done chan bool) {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r%s... Done!     \n", message)
			return
		default:
			fmt.Printf("\r%s %s", spinChars[i], message)
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinChars)
		}
	}
}

func callOpenAI(apiKey string, prompt string, endpoint string, model string) (*gw_model.Config, openai.ChatCompletionResponse, error) {
	var client *openai.Client

	// Create client with custom endpoint if provided
	if endpoint != "" {
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = endpoint
		client = openai.NewClientWithConfig(config)
	} else {
		client = openai.NewClient(apiKey)
	}

	// Create a channel to control the spinner
	done := make(chan bool)
	go startSpinner("Thinking. The process can take a few minutes to finish", done)

	resp, err := client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model:           model,
			Messages:        []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
			ReasoningEffort: "high",
		},
	)

	// Stop the spinner
	done <- true

	if err != nil {
		return nil, openai.ChatCompletionResponse{}, xerrors.Errorf("fail to call open-ai: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"Total tokens":  resp.Usage.TotalTokens,
		"Input tokens":  resp.Usage.PromptTokens,
		"Output tokens": resp.Usage.CompletionTokens,
	}).Info("OpenAI usage:")

	answerText := strings.TrimSpace(resp.Choices[0].Message.Content)
	if err := saveToFile("open-ai-raw.log", answerText); err != nil {
		return nil, openai.ChatCompletionResponse{}, xerrors.Errorf("unable to save raw response: %w", err)
	}

	var res gw_model.Config
	if err := yaml.Unmarshal([]byte(answerText), &res); err != nil {
		return nil, openai.ChatCompletionResponse{}, xerrors.Errorf("unable to unmarshal response: %w", err)
	}
	return &res, resp, nil
}

// Get schema from database config if it exists
func getSchemaFromConfig(databaseType string, configRaw []byte) string {
	// Default schema is empty
	schema := ""

	// Try to parse the config to get the schema for any database type
	var generalConfig struct {
		Schema string `yaml:"schema"`
	}

	if err := yaml.Unmarshal(configRaw, &generalConfig); err == nil {
		if generalConfig.Schema != "" {
			schema = generalConfig.Schema
		}
	}

	// Handle special case for PostgreSQL where public is the default schema
	if databaseType == "postgres" && schema == "" {
		schema = "public"
	}

	return schema
}
