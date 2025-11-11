package migrations

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("v0.4.3")
		tracer.Printf("go ...")

		// update collection `certificate`
		//   - rename field `acmeRenewed` to `isRenewed`
		//   - add field `isRevoked`
		//   - add field `validityInterval`
		{
			collection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(11, []byte(`{
				"hidden": false,
				"id": "number2453290051",
				"max": null,
				"min": null,
				"name": "validityInterval",
				"onlyInt": false,
				"presentable": false,
				"required": false,
				"system": false,
				"type": "number"
			}`)); err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(14, []byte(`{
				"hidden": false,
				"id": "bool810050391",
				"name": "isRenewed",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "bool"
			}`)); err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(15, []byte(`{
				"hidden": false,
				"id": "bool3680845581",
				"name": "isRevoked",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "bool"
			}`)); err != nil {
				return err
			}

			if err := app.Save(collection); err != nil {
				return err
			}

			if _, err := app.DB().NewQuery("UPDATE certificate SET validityInterval = (STRFTIME('%s', validityNotAfter) - STRFTIME('%s', validityNotBefore))").Execute(); err != nil {
				return err
			}
		}

		// adapt to new workflow data structure
		{
			type dWorkflowNode struct {
				Id     string           `json:"id"`
				Type   string           `json:"type"`
				Data   map[string]any   `json:"data"`
				Blocks []*dWorkflowNode `json:"blocks,omitempty,omitzero"`
			}

			var deepMigrateNode func(node *dWorkflowNode) (_node *dWorkflowNode, _migrated bool)
			var deepMigrateNodes func(nodes []*dWorkflowNode) (_nodes []*dWorkflowNode, _migrated bool)
			deepMigrateNode = func(node *dWorkflowNode) (*dWorkflowNode, bool) {
				migrated := false

				if node.Type == "bizApply" {
					if node.Data != nil {
						if _, ok := node.Data["config"]; ok {
							nodeCfg := node.Data["config"].(map[string]any)
							if nodeCfg["keySource"] == nil || nodeCfg["keySource"] == "" {
								nodeCfg["keySource"] = "auto"
								node.Data["config"] = nodeCfg
								migrated = true
							}
						}
					}
				} else if node.Type == "bizUpload" {
					if node.Data != nil {
						if _, ok := node.Data["config"]; ok {
							nodeCfg := node.Data["config"].(map[string]any)
							if nodeCfg["source"] == nil || nodeCfg["source"] == "" {
								nodeCfg["source"] = "form"
								node.Data["config"] = nodeCfg
								migrated = true
							}
						}
					}
				}

				if len(node.Blocks) > 0 {
					if newBlocks, changed := deepMigrateNodes(node.Blocks); changed {
						node.Blocks = newBlocks
						migrated = true
					}
				}

				return node, migrated
			}
			deepMigrateNodes = func(nodes []*dWorkflowNode) ([]*dWorkflowNode, bool) {
				migrated := false

				for i, node := range nodes {
					if newNode, changed := deepMigrateNode(node); changed {
						nodes[i] = newNode
						migrated = true
					}
				}

				return nodes, migrated
			}

			// update collection `workflow`
			//   - migrate field `graphDraft` / `graphContent`
			{
				collection, err := app.FindCollectionByNameOrId("tovyif5ax6j62ur")
				if err != nil {
					return err
				}

				records, err := app.FindAllRecords(collection)
				if err != nil {
					return err
				}

				for _, record := range records {
					changed := false

					graphDraft := make(map[string]any)
					if err := record.UnmarshalJSONField("graphDraft", &graphDraft); err == nil {
						if _, ok := graphDraft["nodes"]; ok {
							nodes := make([]*dWorkflowNode, 0)
							if err := mapstructure.Decode(graphDraft["nodes"], &nodes); err != nil {
								return err
							}

							if newNodes, migrated := deepMigrateNodes(nodes); migrated {
								graphDraft["nodes"] = newNodes
								record.Set("graphDraft", graphDraft)
								changed = true
							}
						}
					}

					graphContent := make(map[string]any)
					if err := record.UnmarshalJSONField("graphContent", &graphContent); err == nil {
						if _, ok := graphContent["nodes"]; ok {
							nodes := make([]*dWorkflowNode, 0)
							if err := mapstructure.Decode(graphContent["nodes"], &nodes); err != nil {
								return err
							}

							if newNodes, migrated := deepMigrateNodes(nodes); migrated {
								graphContent["nodes"] = newNodes
								record.Set("graphContent", graphContent)
								changed = true
							}
						}
					}

					if changed {
						if err := app.Save(record); err != nil {
							return err
						}

						tracer.Printf("record #%s in collection '%s' updated", record.Id, collection.Name)
					}
				}
			}

			// update collection `workflow_run`
			//   - migrate field `graph`
			{
				collection, err := app.FindCollectionByNameOrId("qjp8lygssgwyqyz")
				if err != nil {
					return err
				}

				records, err := app.FindAllRecords(collection)
				if err != nil {
					return err
				}

				for _, record := range records {
					changed := false

					graph := make(map[string]any)
					if err := record.UnmarshalJSONField("graph", &graph); err == nil {
						if _, ok := graph["nodes"]; ok {
							nodes := make([]*dWorkflowNode, 0)
							if err := mapstructure.Decode(graph["nodes"], &nodes); err != nil {
								return err
							}

							if newNodes, migrated := deepMigrateNodes(nodes); migrated {
								graph["nodes"] = newNodes
								record.Set("graph", graph)
								changed = true
							}
						}
					}

					if changed {
						if err := app.Save(record); err != nil {
							return err
						}

						tracer.Printf("record #%s in collection '%s' updated", record.Id, collection.Name)
					}
				}
			}
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}
