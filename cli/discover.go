package cli

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/centralmind/gateway/logger"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/centralmind/gateway/connectors"
	gw_model "github.com/centralmind/gateway/model"
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
	var llmLogFile string

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

			logrus.Info("Step 4: Prepare the prompt for the AI")
			discoverPrompt := generateDiscoverPrompt(databaseType, extraPrompt, tablesToGenerate, getSchemaFromConfig(databaseType, configRaw))
			if err := saveToFile(promptFile, discoverPrompt); err != nil {
				logrus.Error("failed to save prompt:", err)
			}

			logrus.Debugf("Prompt saved locally to %s", promptFile)
			logrus.Info("✅ Step 4 completed. Done.")
			logrus.Info("\r\n")

			// Call API
			logrus.Info("Step 5: Use AI to design the API")
			response, err := makeDiscoverQuery(DiscoverQueryParams{
				LLMLogFile:    llmLogFile,
				Provider:      aiProvider,
				Endpoint:      aiEndpoint,
				APIKey:        aiAPIKey,
				Model:         aiModel,
				MaxTokens:     aiMaxTokens,
				Temperature:   aiTemperature,
				Reasoning:     aiReasoning,
				BedrockRegion: bedrockRegion,
				VertexRegion:  vertexAIRegion,
				VertexProject: vertexAIProject,
			}, discoverPrompt)

			if err != nil {
				logrus.Error("Failed to call the LLM:", err)
				return err
			}

			config := response.Config

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

			logrus.Info("✅ All steps completed. Done.")
			logrus.Info("\r\n")
			logrus.Info("--- Execution Statistics ---")
			logrus.Infof("Total time taken: "+yellow+"%v"+reset, duration.Round(time.Second))
			logrus.Infof("Tokens used: "+yellow+"%d"+reset+" (Estimated cost: "+violet+"$%.4f"+reset+")",
				response.Conversation.Usage.TotalTokens, response.CostEstimate)
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

			logrus.Infof("Total number of columns containing PII data: "+yellow+"%d"+reset, piiColumnsCount)

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
	cmd.Flags().StringVar(&aiAPIKey, "ai-api-key", "", "AI API token")
	cmd.Flags().StringVar(&bedrockRegion, "bedrock-region", "", "Bedrock region")
	cmd.Flags().StringVar(&vertexAIRegion, "vertexai-region", "", "Vertex AI region")
	cmd.Flags().StringVar(&vertexAIProject, "vertexai-project", "", "Vertex AI project")
	cmd.Flags().StringVar(&aiModel, "ai-model", "", "AI model to use")
	cmd.Flags().IntVar(&aiMaxTokens, "ai-max-tokens", 0, "Maximum tokens to use")
	cmd.Flags().Float32Var(&aiTemperature, "ai-temperature", -1.0, "AI temperature")
	cmd.Flags().BoolVar(&aiReasoning, "ai-reasoning", true, "Enable reasoning")

	cmd.Flags().StringVar(&output, "output", "gateway.yaml", "Resulted YAML path")
	cmd.Flags().StringVar(&extraPrompt, "prompt", "generate reasonable set of APIs for this data", "Custom input to generate APIs")
	cmd.Flags().StringVar(&promptFile, "prompt-file", filepath.Join(logger.DefaultLogDir(), "prompt_default.txt"), "Path to save the generated prompt")
	cmd.Flags().StringVar(&llmLogFile, "llm-log", filepath.Join(logger.DefaultLogDir(), "llm_raw_response.log"), "Path to save the raw LLM response")

	return cmd
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
