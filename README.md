# DDD - Data Definition Data

DDD is a format for describing what data should look like.

# Data Types

There are 4 types of data in DDD:

1. List (list)
2. Key Dict (keydict)
3. Row Dict (rowdict)
4. Value Requirement

Each of these is wrapped in a dictionary ({}) so that they can contain arguments.  For example:

**List:**

```
{"list": []}
```

Lists are sequences of things.  A list can hold any type of DDD data: list, keydict, rowdict, value requirement.

**Key Dict:**

```
{"keydict": {}}
```

Key Dicts hold some set of given keys, which are specified in their requirements by the keys.  The keys are typically
specificied as a value (ex: string), and can be any DDD value: list, keydict or value requirement.  (Row Dicts only appear in lists)

**Row Dict:**

```
{"rowdict": {}}
```

Row Dicts are meant to repeat, and so only appear inside lists.  Their requirements are switched on a single field,
such as "type", which allows many different row formats to exist, for the dictionaries in the list.

**Value Requirement:**

```
{}
```

Value Requirements appear in lists, or on the Value side of Key Dict or Row Dict Key/Value elements.


# Examples

## JSON Format: Pairs of strings, that can be repeated in the 2nd and 3rd list depths

```
[
  [
    [
      "__get.param.udn.__execute",
      ""
    ]
  ]
]
```

### DDD

```{
     "list": [
       {
         "list": [
           {
             "list": [
               {
                 "type": "string"
               },
               {
                 "type": "string"
               }
             ]
           },
           "..."
         ]
       },
       "..."
     ]
   }
```

## JSON Format: Description of defaults and a table name

```
{
  "table": "web_site_page",
  "defaults": {
    "web_site_id": 1,
    "base_page_web_site_page_widget_id": 46,
    "login_required": true
  }
}
```

### DDD

```
{
  "keydict": {
    "table": {
      "type": "string"
    },
    "defaults": {
      "keydict": {
        "*": {
          "type": "any"
        }
      }
    }
  }
}
```

## JSON Format: Description of form field layout

```
{
  "dialog_title": "Edit Stored UDN Function",
  "form":[
    [
      {"name":"name", "label":"Name", "placeholder":"widget name", "icon":"icon-file-text", "type":"text", "size":"6", "value":"", "info":"", "color":""},
      {"name":"udn_stored_function_domain_id", "label":"Widget Type", "placeholder":"", "icon":"icon-city", "type":"select", "size":"6", "value":"", "info":"", "color":"", "value_match":"select_option_match", "value_nomatch":"select_option_nomatch", "items": "__data_filter.udn_stored_function_domain.{}.__array_map_remap.{name=name,value=_id}"}
    ],
    [
      {"name":"udn_data_json", "label":"Info", "placeholder":"", "icon":"icon-atom", "type":"ace", "size":"12", "value":"", "info":"", "color":"", "format": "json"}
    ],
    [
      {"name":"_id", "type":"hidden", "value":"", "size":"4", "color":"", "label":"ID"},
      {"name":"_web_data_widget_instance_id", "type":"hidden", "value":"", "size":"4", "color":"", "label":"Data Widget"},
      {"name":"_table", "type":"hidden", "value":"", "size":"4", "color":"", "label":"Table"}
    ]
  ]
}
```

### DDD

```
{
    "keydict": {
        "dialog_title": {
            "type": "string"
        },
        "form": {
            "list": [{
                "list": [{
                    "rowdict": {
                        "switch_field": "type",
                        "switch_rows": {
                            "text": {
                                "name": {"type": "string"},
                                "label": {"type": "string"},
                                "placeholder": {"type": "string"},
                                "icon": {"type": "string"},
                                "size": {"type": "int", "min": 2, "max": 12},
                                "value": {"type": "string"}
                            },
                            "select": {
                                "name": {"type": "string"},
                                "label": {"type": "string"},
                                "placeholder": {"type": "string"},
                                "icon": {"type": "string"},
                                "size": {"type": "int", "min": 2, "max": 12},
                                "value": {"type": "string"},
                                "value_match": {"type": "string"},
                                "value_nomatch": {"type": "string"},
                                "null_message": {"type": "string", "optional": true},
                                "items": {"type": "udn"}
                            },
                            "hidden": {
                                "name": {"type": "string"},
                                "value": {"type": "string"}
                            }
                        }
                    }
                }, "..."]
            }, "..."]
        }
    }
}
```