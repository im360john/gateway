package cli

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	- The final output must contain only valid JSON with no additional commentary, explanations, or markdown formatting.
	- The JSON configuration must strictly adhere to the provided JSON schema, including all required fields.
	- Description of API endpoints should have also an example, to help chatbot to use it.
	- All descriptions and summary must not have any sensetive information/data from security point of view including database types, password and etc.
	- All SQL queries must be Pure SQL that will be used in golang SQLx on top of database - {database_type} and be fully parameterized (using named parameters) to prevent SQL injection.
	- All API endpoints must have output schemas.
	- All SQL queries must be verified that they will not return array of data where expected one item.
	- SQL queries should be optimized for {database_type} and use appropirate indexes.
	- Endpoints that return lists must include pagination parameters (offset and limit).
	- Consistent Endpoint Definitions: Each table defined in the DDL should have corresponding endpoints as specified by the JSON schema, including method, path, description, SQL query, and parameters.
	- Sensitive Data Handling: If any columns contain sensitive data, they must be flagged appropriately (e.g., using a "pii" flag).
	- Each Parameter in API endpoints must have default value taken from corresponded example rows
`
)

var (
	//go:embed api_config_schema.json
	apiConfigSchema []byte
)

type TableData struct {
	Columns []gw_model.ColumnSchema
	Name    string
	Sample  []map[string]any
}

func Discover(configPath *string) *cobra.Command {
	var databaseType string
	var tables []string
	var openAPIKey string
	var output string
	var extraPrompt string
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover generates gateway config",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			startTime := time.Now()
			logrus.Infof("Step 1: Read configs")
			configRaw, err := os.ReadFile(*configPath)
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
			logrus.Infof("Step 2: Discover data")
			allTables, err := connector.Discovery(context.Background())
			if err != nil {
				return err
			}
			logrus.Infof("Step 2: Found: %v tables", len(allTables))
			tableSet := map[string]bool{}
			for _, table := range tables {
				tableSet[table] = true
			}
			if len(tables) == 0 {
				for _, table := range allTables {
					tableSet[table.Name] = true
				}
			}
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
					Columns: table.Columns,
					Name:    table.Name,
					Sample:  sample,
				})
			}
			logrus.Info("Step 3: Prepare prompt to AI")
			fullPrompt := generatePrompt(databaseType, extraPrompt, tablesToGenerate)
			promptFilename := "prompt.txt"
			if err := saveToFile(promptFilename, fullPrompt); err != nil {
				logrus.Error("failed to save prompt:", err)
			}

			logrus.Infof("Step 3 done. Prompt: %s", promptFilename)

			logrus.Info("Step 4: Do AI Magic")
			config, err := callOpenAI(openAPIKey, fullPrompt)
			if err != nil {
				logrus.Error("failed to call OpenAI:", err)
				return
			}

			config.Database.Type = databaseType
			config.Database.Connection = string(configRaw)

			configData, err := yaml.Marshal(config)
			if err != nil {
				logrus.Error("yaml failed:", err)
				return
			}

			if err := saveToFile(output, string(configData)); err != nil {
				logrus.Error("failed:", err)
			}

			logrus.Infof("✅ API schema saved в %s", output)

			logrus.Infof("Done: in %v", time.Since(startTime))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&tables, "tables", nil, "List of table to include")
	cmd.Flags().StringVar(&databaseType, "db-type", "postgres", "Type of database")
	cmd.Flags().StringVar(&openAPIKey, "open-ai-key", "open-ai-key", "OpenAI token")
	cmd.Flags().StringVar(&output, "output", "gateway.yaml", "Resulted yaml path")
	cmd.Flags().StringVar(&extraPrompt, "prompt", "generate reasonable set of API-s for this data", "Custom input to generate API-s")
	return cmd
}

func generatePrompt(databaseType, extraPrompt string, tables []TableData) string {
	res := "I need a config for an automatic API that will be used by another AI bot or LLMs..."
	res += "\n"
	res += strings.ReplaceAll(basePrompt, "{database_type}", databaseType)
	res += "\n" + string(apiConfigSchema) + "\n" + extraPrompt + "\n\n"
	for _, table := range tables {
		res += fmt.Sprintf(`
<%[1]s number_columns=%[5]v number_rows=%[4]v>
schema:
%[2]s
---
data_sample:
%[3]s
</%[1]s>

`, table.Name, yamlify(table.Columns), yamlify(table.Sample), len(table.Sample), len(table.Columns))
	}

	return res
}

func yamlify(sample any) string {
	raw, _ := yaml.Marshal(sample)
	return string(raw)
}

func saveToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func callOpenAI(apiKey string, prompt string) (*gw_model.Config, error) {
	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model:           "o3-mini",
			Messages:        []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
			ReasoningEffort: "high",
		},
	)
	if err != nil {
		return nil, errors.Errorf("fail to call open-ai: %w", err)
	}
	logrus.Infof("Step 4: open-ai usage: %v", resp.Usage)

	answerText := strings.TrimSpace(resp.Choices[0].Message.Content)
	if err := saveToFile("open-ai-raw.log", answerText); err != nil {
		return nil, errors.Errorf("unable to save raw response: %w", err)
	}

	var res gw_model.Config
	if err := yaml.Unmarshal([]byte(answerText), &res); err != nil {
		return nil, errors.Errorf("unable to unmarshal response: %w", err)
	}
	return &res, nil
}
