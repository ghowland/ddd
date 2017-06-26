# DDD - Data Definition Data

**DDD is a format for describing what data should look like.**

It can be used for validation as well as population.  It was created to guide population of "arbitrary data", which has a format, but there are many different formats the users must populate correctly, so a guide was needed to ensure the creation of correct data.

*NOTE: DDD is not meant to be able to model any possible data, but it can model many formats of data with very simple descriptions.  I haven't tested it to see what it cannot create/validate yet, but for all the things I currently need to validate/populate, it works well.  DDD probably can't describe XML/HTML, and I'm not sure if it can describe itself; neither of which are goals for the format.*

# Data Types

There are 4 types of data in DDD:

1. List (`list`)
2. Key Dict (`keydict`)
3. Row Dict (`rowdict`)
4. Value Requirement (`type`)

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
{"type": "string"}
```

Value Requirements appear in lists, or on the Value side of Key Dict or Row Dict Key/Value elements.

### Optionality

```
"optional": false
```

The above field is a default value in all of our types.  Any of them can set this to true, and then do not need to exist.  If it is not specified, one of the specified requirements must be met in the data to be valid.

### Value Requirements Constraints

```
min: 0,
max: 100
```

Value Requirements can have different constraints added to them as well, such as min/max for numbers, which are
also used for strings as min/max size.  Other constraints can also be added for Value Requirements constraining in the implementation.


# Examples

## JSON Format: Pairs of strings in lists, that can be repeated at the 2nd and 3rd list depths

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

In the above format, the inner 2 string pairs are fixed in length, but we can have any number of those 2-tuples of strings, or any number of the next higher level of list element (which contain the groupings of 2-tuple strings)

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

3 embedded lists, where the inner 2 lists are variadic, which means they can take N lists of 2-tuple strings, or N of the next higher level list containers.

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

This data has a fixed field (`table`), and then a dictionary that can have any number of fields (`defaults`) with any values.

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

We use a `"keydict"` at the top level, because we are basing this on dictionary keys, and we need to be in a list to use Row Dicts.

Defaults uses the `"*"` character to specify any field name can go here (any number of them as well), and the value requirements are type `"any"`, so any type of data can be stored here.

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
      {"name":"_id", "type":"hidden", "value":""},
      {"name":"_web_data_widget_instance_id", "type":"hidden", "value":""},
      {"name":"_table", "type":"hidden", "value":""}
    ]
  ]
}
```

This is similar to the previous example, with a fixed string in `"dialog_title"`, but the `"form"` is a list of lists (rows and columns) for fields in an Edit Form layout.

Because we are using rows, and we want to have repeating dictionaries in those rows with correct field values/types and the correct names/count of fields, we use a Row Dict, and perform a "switch" for the record type on the `"type"` field, which is the key for this Edit Form field layout data.

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
                                "items": {"type": "string"}
                            },
                            "hidden": {
                                "name": {"type": "string"},
                                "value": {"type": "string"}
                            }
                        }
                    }, "..."]
                }, "..."]
            }
        }
    }
}
```

`"form"` gets a list (columns), in a list (rows), and inside there specifies a Row Dict, which is variadic, so that the columns can have more than one of these dictionaries per column-list.

Because the first `"rowdict"` entry doesnt contain a `"optional": true` field, there must be at least 1 dictionary per column-list.  If this contained `"optional": false`, then the column-lists could be empty.

We set `"switch_field"` to `"type"`, so we will be using the `"type"` field in the dictionaries to set which possible fields they can have.  The available values for the `"type"` field are the keys in `"switch_rows"`, so type can be: `text`, `select` and `hidden`

The `"rowdict"` in the columns-list is followed by a variadic `"..."`, so there can be more than 1 dictionary in the column-list, which makes more than 1 column.

The column-list `"list"` is followed by a variadic `"..."`, so there can be multiple rows as the column-list is repeated.

The outer list contains the list of rows.

## Just a number

```
6
```

If we want my data to be a single value.

### DDD

```
{"type": "int", "optional": false}
```

This ensures an integer must exist in this data.  If we make it optional, it could be an empty string or return None/nil/etc, depending on the implementation.

## JSON list of things

```
["Bob", 5.1, [1,2,3,1000]]
```

### DDD

```
{"list":
  [
    {"type": "string"},
    {"type": "float"},
    {"list":
      [
        {"type": "int", "optional": true, "min": 0, "max": 1001},
        "..."
      ]
    }
  ]
}
```

This is a list with 3 fixed elements:  a string, a float, and a list.

The inner list contains 0 to an infinite number of integers, with a minimum value of 0, and a maximum value of 1001.

