package ddd

import (
	"strings"
	"fmt"
	"strconv"
	"sort"
	. "yudien/yudienutil"
)


const (
	type_int				= iota
	type_float				= iota
	type_string				= iota
	type_string_force		= iota	// This forces it to a string, even if it will be ugly, will print the type of the non-string data too.  Testing this to see if splitting these into 2 yields better results.
	type_array				= iota	// []interface{} - takes: lists, arrays, maps (key/value tuple array, strings (single element array), ints (single), floats (single)
	type_map				= iota	// map[string]interface{}
)

/*

func DddSet(position_location string, data_location string, save_data map[string]interface{}, ddd_id int, udn_data map[string]interface{}) {
	ddd := DatamanGet("ddd", ddd_id)

	// Get our positional info
	position_info := _DddGetPositionInfo(position_location, udn_data)

}

func DddValidate(data_location string, ddd_id int, udn_data map[string]interface{}) []map[string]interface{} {
	ddd := DatamanGet("ddd", ddd_id)

	result := make([]map[string]interface{}, 0)
	return result
}

func DddDelete(position_location string, data_location string, ddd_id int, udn_data map[string]interface{}) {
	ddd := DatamanGet("ddd", ddd_id)

	// Get our positional info
	position_info := _DddGetPositionInfo(position_location, udn_data)

}
*/

/*
func _DddGetPositionInfo(position_location string, udn_data map[string]interface{}) map[string]interface{} {
	current_info := MapGet(MakeArray(position_location), udn_data).(map[string]interface{})

	// If current_info is not set up properly, set it up
	if current_info == nil || current_info["x"] == nil {
		current_info = make(map[string]interface{})

		current_info["location"] = "0"
		current_info["x"] = 0
		current_info["y"] = 0
	}

	return current_info
}
*/

func DddMove(position_location string, move_x int64, move_y int64) string {
	//NOTE(g): This function doesnt check if the new position is valid, that is done by DddGet() which returns the DDD info at the current position (if valid, or nil

	parts := strings.Split(position_location, ".")

	fmt.Printf("DDD Move: Parts: %v\n", parts)

	// Only allow X or Y movement, not both.  This isnt a video game.
	if move_x != 0 {
		if move_x == 1 {
			fmt.Printf("DDD Move: RIGHT\n")
			// Moving to the right, we just add a .0 to the current location
			return fmt.Sprintf("%s.0", position_location)
		} else {
			fmt.Printf("DDD Move: LEFT\n")
			if len(parts) > 1 {
				parts = parts[0:len(parts)-1]
				return strings.Join(parts, ".")
			} else {
				fmt.Printf("DDD Move: Cant move left\n")
				// Else, we only have 1 location part, we cant reduce this, so return the initial location
				return position_location
			}
		}
	} else if move_y != 0 {
		last_part := parts[len(parts)-1]
		last_part_int, _ := strconv.Atoi(last_part)

		if move_y == 1 {
			fmt.Printf("DDD Move: DOWN\n")
			// Moving down, increment the last_part_int
			last_part_int++
			parts[len(parts)-1] = strconv.Itoa(last_part_int)

			return strings.Join(parts, ".")
		} else {
			fmt.Printf("DDD Move: UP\n")
			// Moving up, decrement the last_part_int
			last_part_int--
			//if last_part_int < 0 {
			//	last_part_int = 0
			//}
			parts[len(parts)-1] = strconv.Itoa(last_part_int)

			return strings.Join(parts, ".")
		}
	}

	fmt.Printf("DDD Move: No Change\n")

	// No change in position, return the same string we received
	return position_location
}

func DddGet(position_location string, data_location string, ddd_data map[string]interface{}, udn_data map[string]interface{}) interface{} {
	// Get the DDD Node that describes this position
	//ddd_node := DddGetNode(position_location, ddd_data, udn_data)

	//TODO(g): SECOND!    We know the DDD information, so we navigate the same way we did DDD, but we get the data
	//
	//	What if it isnt available?  We return an error.  How?
	//
	//	??	How		??
	//		???
	//
	// Copy the looping code into all the functions, dont worry about generalizing initially, just get it working.
	//


	result := 1
	return result
}


func _DddGetNodeCurrent(cur_data map[string]interface{}, cur_record_data interface{}, cur_pos int, processed_parts []int, cur_parts []string) (string, map[string]interface{}, interface{}) {
	if cur_data["keydict"] != nil {
		// The cur_pos will be selected based on the sorted values, because they are map-keys, they are out of order.  Once sorted, they are accessed as an array index

		keys := MapKeys(cur_data["keydict"].(map[string]interface{}))

		fmt.Printf("DddGetNodeCurrent: keydict: Keys: %v\n", keys)

		// We didnt find it, so return nil
		if cur_pos >= len(keys) || cur_pos < 0 {
			return "nil", nil, nil
		}

		selected_key := keys[cur_pos]

		fmt.Printf("DddGetNodeCurrent: keydict: Selected Key: %s\n", selected_key)

		result_cur_data := cur_data["keydict"].(map[string]interface{})[selected_key].(map[string]interface{})

		cur_record_data_map := GetResult(cur_record_data, type_map).(map[string]interface{})

		result_cur_record_data := make(map[string]interface{})
		if cur_record_data_map[selected_key] != nil {
			result_cur_record_data = GetResult(cur_record_data_map[selected_key], type_map).(map[string]interface{})
		}

		return fmt.Sprintf("Key: %s", selected_key), result_cur_data, result_cur_record_data

	} else if cur_data["rowdict"] != nil {
		// The rowdict is inside a list, but must be further selected based on the selection field, which will determine the node
		//TODO(g): ...
		return "RowDict", cur_data, cur_record_data

	} else if cur_data["list"] != nil {
		fmt.Printf("DDDGET:LIST: %T\n", cur_data["list"])
		cur_data_list := cur_data["list"].([]interface{})

		// Using the cur_pos as the index offset, this works up until the "variadic" node (if present)
		if cur_pos >= 0 && cur_pos < len(cur_data_list) {
			result_cur_data := cur_data_list[cur_pos].(map[string]interface{})

			var result_cur_record_data interface{}

			cur_record_data_array := GetResult(cur_record_data, type_array).([]interface{})

			if len(cur_record_data_array) > cur_pos {
				result_cur_record_data = cur_record_data_array[cur_pos]
			} else {
				result_cur_record_data = nil
			}

			return fmt.Sprintf("Index: %d", cur_pos), result_cur_data, result_cur_record_data
		} else {
			return "nil", nil, nil
		}

	} else if cur_data["type"] != nil {
		// This is a raw data node, and should not have any indexing, only "0" for it's location position
		if cur_pos == 0 {
			return "TBD: Get Label", cur_data, cur_record_data
		} else {
			return "nil", nil, nil
		}

	} else if cur_data["variadic"] != nil {
		// I think I have to backtrack to a previous node then?  Parent node?
		if cur_pos == 0 {
			return fmt.Sprintf("Variadic: %d", cur_pos), cur_data, cur_record_data
		} else {
			return "nil", nil, nil
		}

	} else {
		//TODO(g): Replace this panic with a non-fatal error...  But the DDD is bad, so report it?
		//panic(fmt.Sprintf("Unknown DDD node: %v", cur_data))
		return "nil", nil, nil
	}

	return "Unknown", cur_data, cur_record_data
}

func DddGetNode(position_location string, ddd_data map[string]interface{}, data_record interface{}, udn_data map[string]interface{}) (string, map[string]interface{}, interface{}) {
	cur_parts := strings.Split(position_location, ".")
	cur_label := ""
	fmt.Printf("DDD Get Node: Parts: %s: %v\n", position_location, cur_parts)

	// Current position starts from ddd_data, and then we navigate it, and return it when we find the node
	cur_data := ddd_data
	cur_record_data := data_record

	processed_parts := make([]int, 0)

	// The first "0" is always "0", and is the base cur_data, so let's pop it off
	if len(cur_parts) > 1 {
		// Add the part we just processed to our processed_parts slice to keep track of them
		cur_pos, _ := strconv.Atoi(cur_parts[0])
		processed_parts = append(processed_parts, cur_pos)

		fmt.Printf("DddGetNode: Removing first part: %v\n", cur_parts)
		cur_parts = cur_parts[1:len(cur_parts)]
		fmt.Printf("DddGetNode: Removed first part: %v\n", cur_parts)
	} else {
		if position_location == "0" {
			// There are no other parts, so we have the data
			fmt.Printf("DddGetNode: First part is '0': %s\n", position_location)
			return "The Beginninging", cur_data, cur_record_data
		} else {
			// Asking for data which cannot exist.  The first part can only be 0
			fmt.Printf("DddGetNode: First part is only part, and isnt '0': %s\n", position_location)
			return "The Somethingelseinging", nil, nil
		}
	}



	// As long as we still have cur_parts, keep going.  If we dont return in this block, we will have an empty result
	for len(cur_parts) > 0 {
		cur_pos, _ := strconv.Atoi(cur_parts[0])
		fmt.Printf("DDD Move: Step: Parts: %v   Current: %d  Cur Node: %s  Cursor Data: %s\n", cur_parts, cur_pos, SnippetData(cur_data, 80), SnippetData(cur_record_data, 80))

		cur_label, cur_data, cur_record_data = _DddGetNodeCurrent(cur_data, cur_record_data, cur_pos, processed_parts, cur_parts)

		// Add the part we just processed to our processed_parts slice to keep track of them
		processed_parts = append(processed_parts, cur_pos)

		// Pop off the first element, so we keep going
		if len(cur_parts) > 1 {
			cur_parts = cur_parts[1:len(cur_parts)]
		} else {
			cur_parts = make([]string, 0)
		}

		// If we have nothing left to process, return the result
		if len(cur_parts) == 0 {
			fmt.Printf("DddGetNode: Result: %s: Node Data: %s  Cursor Data: %s\n", position_location, SnippetData(cur_data, 80), SnippetData(cur_record_data, 80))
			return cur_label, cur_data, cur_record_data
		} else if cur_data["type"] != nil || cur_data["variadic"] != nil || cur_data["rowdict"] != nil {
			return cur_label, nil, nil
		}
	}

	// No data at this location, or we would have returned it already
	fmt.Printf("DddGetNode: No result, returning nil: %v\n", cur_parts)
	return "nil", nil, nil
}

func GetDddNodeSummary(cur_label string, cur_data map[string]interface{}) string {
	// This is our result, setting to unknown, which should never be displayed
	summary := "Unknown: FIX"

	if cur_data["keydict"] != nil {
		keys := MapKeys(cur_data["keydict"].(map[string]interface{}))

		summary = fmt.Sprintf("%s: KeyDict: %v", cur_label, strings.Join(keys, ", "))

	} else if cur_data["rowdict"] != nil {
		keys := MapKeys(cur_data["rowdict"].(map[string]interface{})["switch_rows"].(map[string]interface{}))

		summary = fmt.Sprintf("%s: RowDict: Rows: %d:  %v", cur_label, len(cur_data["rowdict"].(map[string]interface{})), strings.Join(keys, ", "))

	} else if cur_data["list"] != nil {
		cur_list := cur_data["list"].([]interface{})

		item_summary := make([]string, 0)
		for _, item := range cur_data["list"].([]interface{}) {
			item_summary = append(item_summary, GetDddNodeSummary("", item.(map[string]interface{})))
		}
		item_summary_str := strings.Join(item_summary, ", ")


		summary = fmt.Sprintf("%s: List (%d): %s", cur_label, len(cur_list), item_summary_str)

	} else if cur_data["type"] != nil {
		summary = fmt.Sprintf("%s: Data Item: Type: %s", cur_label, cur_data["type"])

	} else if cur_data["variadic"] != nil {
		summary = fmt.Sprintf("%s: Variadic", cur_label)
	}

	// Crop long summaries
	if len(summary) > 60 {
		summary = summary[0:60]
	}

	return summary
}



func GetFieldMapFromSpec(data map[string]interface{}, label string, name string) map[string]interface{} {
	field_map := make(map[string]interface{})

	if data["type"] == "string" || data["type"] == "int" || data["type"] == "boolean" {
		icon := "icon-make-group"
		if data["icon"] != nil {
			icon = data["icon"].(string)
		}

		size := 12
		if data["size"] != nil {
			size = int(data["size"].(float64))
		}

		field_map = map[string]interface{}{
			"color": "primary",
			"icon": icon,
			"info": "",
			"label": label,
			"name": name,
			"placeholder": "",
			"size": size,
			"type": "text",
			"value": "",
		}
	}

	return field_map
}

func DddRenderNode(position_location string, ddd_id int64, temp_id int64, ddd_label string, ddd_node map[string]interface{}, ddd_cursor_data interface{}) []interface{} {
	rows := make([]interface{}, 0)

	//// Add the current row, so we work with them
	//cur_row := make([]interface{}, 0)
	//rows = append(rows, cur_row)

	if ddd_node["type"] != nil {
		field_name := fmt.Sprintf("ddd_node_%s", position_location)
		new_html_field := GetFieldMapFromSpec(ddd_node, ddd_label, field_name)
		rows = AppendArray(rows, new_html_field)
	} else if ddd_node["keydict"] != nil {
		html_element_name := fmt.Sprintf("ddd_node_%s", position_location)

		// Keydict select fields, navs to them, so we dont have to button nav
		new_html_field := map[string]interface{}{
			"color": "primary",
			"icon": "icon-make-group",
			"info": "",
			"label": ddd_label,
			"name": html_element_name,
			"placeholder": "",
			"size": "12",
			"type": "select",
			"value": "",
			"value_match":"select_option_match",
			"value_nomatch":"select_option_nomatch",
			"null_message": "- Select to Navigate -",
			"items": fmt.Sprintf("__input.%s", MapKeysToUdnMapForHtmlSelect(position_location, ddd_node["keydict"].(map[string]interface{}))),
			"onchange": fmt.Sprintf("$(this).closest('.ui-dialog-content').dialog('close'); RPC('/api/dwi_render_ddd', {'move_x': 0, 'move_y': 0, 'position_location': $(this).val(), 'ddd_id': %d, 'is_delete': 0, 'web_data_widget_instance_id': '{{{_id}}}', 'web_widget_instance_id': '{{{web_widget_instance_id}}}', '_web_data_widget_instance_id': 34, 'dom_target_id':'dialog_target', 'temp_id': %d})", ddd_id, temp_id),
		}
		rows = AppendArray(rows, new_html_field)
	} else if ddd_node["list"] != nil {
		map_values := make([]string, 0)

		for index, data := range ddd_node["list"].([]interface{}) {
			summary := GetDddNodeSummary(ddd_label, data.(map[string]interface{}))

			new_position := fmt.Sprintf("%s.%d", position_location, index)

			map_values = append(map_values, fmt.Sprintf("{name='%s',value='%s'}", summary, new_position))
		}

		map_value_str := strings.Join(map_values, ",")

		udn_final := fmt.Sprintf("[%s]", map_value_str)

		html_element_name := fmt.Sprintf("ddd_node_%s", position_location)

		// Keydict select fields, navs to them, so we dont have to button nav
		new_html_field := map[string]interface{}{
			"color": "primary",
			"icon": "icon-make-group",
			"info": "",
			"label": ddd_label,
			"name": html_element_name,
			"placeholder": "",
			"size": "12",
			"type": "select",
			"value": "",
			"value_match":"select_option_match",
			"value_nomatch":"select_option_nomatch",
			"null_message": "- Select to Navigate -",
			"items": fmt.Sprintf("__input.%s", udn_final),
			"onchange": fmt.Sprintf("$(this).closest('.ui-dialog-content').dialog('close'); RPC('/api/dwi_render_ddd', {'move_x': 0, 'move_y': 0, 'position_location': $(this).val(), 'ddd_id': %d, 'is_delete': 0, 'web_data_widget_instance_id': '{{{_id}}}', 'web_widget_instance_id': '{{{web_widget_instance_id}}}', '_web_data_widget_instance_id': 34, 'dom_target_id':'dialog_target', 'temp_id': %d})", ddd_id, temp_id),
		}
		rows = AppendArray(rows, new_html_field)
	} else if ddd_node["rowdict"] != nil {
		// Sort by rows and columns, if available, if not, sort them and put them at the end, 1 per row
		unsorted := make([]map[string]interface{}, 0)

		layout := make(map[int]map[int]map[string]interface{})

		//TODO(g): We will assume data initially, so we can start up
		data_switch_field := "text"

		// Select the spec from the switch_field
		selected_row_dict_spec := ddd_node["rowdict"].(map[string]interface{})["switch_rows"].(map[string]interface{})[data_switch_field].(map[string]interface{})

		for key, value := range selected_row_dict_spec {
			value_map := value.(map[string]interface{})

			new_item := make(map[string]interface{})
			new_item[key] = value

			if value_map["x"] != nil && value_map["y"] != nil {
				// Put them in Y first, because we care about ordering by rows first, then columns once in a specific row
				if layout[int(value_map["y"].(float64))] == nil {
					layout[int(value_map["y"].(float64))] = make(map[int]map[string]interface{})
				}
				layout[int(value_map["y"].(float64))][int(value_map["x"].(float64))] = new_item
			} else {
				unsorted = append(unsorted, new_item)
			}
		}

		fmt.Printf("DddRenderNode: RowDict: Layout: %s\n\n", JsonDump(layout))

		// Get the Y keys
		y_keys := make([]int, len(layout))
		i := 0
		for key := range layout {
			y_keys[i] = key
			i++
		}
		sort.Ints(y_keys)

		// Loop over our rows
		max_y := ArrayIntMax(y_keys)
		for cur_y := 0 ; cur_y <= max_y ; cur_y++ {
			//fmt.Printf("DddRenderNode: RowDict: Y: %d\n", cur_y)

			if layout[cur_y] != nil {
				layout_row := layout[cur_y]

				// Get the Y keys
				x_keys := make([]int, len(layout_row))
				i := 0
				for key := range layout_row {
					x_keys[i] = key
					i++
				}
				sort.Ints(x_keys)

				// Loop over our cols
				max_x := ArrayIntMax(x_keys)
				for cur_x := 0 ; cur_x <= max_x ; cur_x++ {
					//fmt.Printf("DddRenderNode: RowDict: Y: %d  X: %d\n", cur_y, cur_x)

					if layout_row[cur_x] != nil {
						layout_item := layout_row[cur_x]

						field_name := fmt.Sprintf("ddd_node_%s__%d_%d", position_location, cur_x, cur_y)

						for layout_key, layout_data := range layout_item {
							layout_data_map := layout_data.(map[string]interface{})
							new_html_field := GetFieldMapFromSpec(layout_data_map, layout_data_map["label"].(string), field_name)

							fmt.Printf("DddRenderNode: RowDict: Y: %d  X: %d:  %s: %v\n", cur_y, cur_x, layout_key, layout_data_map)
							fmt.Printf("%s\n", JsonDump(new_html_field))

							rows = AppendArray(rows, new_html_field)
						}
					} else {
						//TODO(g): Put in empty columns, once I figure it out.  Empty HTML?  That could work...
						fmt.Printf("DddRenderNode: RowDict: Y: %d  X: %d:  No data here, missing column\n", cur_y, cur_x)
					}
				}
			} else {
				//TODO(g): Skip rows?  Or just ignore them?  I dont think I need to render empty rows...
			}

			//// Add the current row, so we work with them
			//cur_row := make([]interface{}, 0)
			//rows = append(rows, cur_row)
		}

		if len(unsorted) > 0 {
			// Sort them based on their ["name"] field
			//TODO(g): TBD...  Then render them after the others in rows
		}
	}

	return rows
}
