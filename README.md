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

For example, here are two different formats of row dictionaries, which are differentiated by their "type" field:

*"type" = "hidden":*

```
{"name":"_id", "type":"hidden", "value":""}
```

*"type" = "text":*

```
{"name":"name", "label":"Name", "placeholder":"widget name", "icon":"icon-file-text", "type":"text", "size":"6", "value":"", "info":"", "color":""}
```

The first one only has 2 fields besides the "type" field (`hidden`), but the second one (`text`) has 8 other fields.  By allowing selection of the field requirements to switched against the "type" value, we can have differently formatted dictionaries in the same list, following these validation rules.

**Value Requirement:**

```
{"type": "string"}
```

Value Requirements appear in lists, or on the Value side of Key Dict or Row Dict Key/Value elements.

# Dynamics and Constraints

### Variadic Elements

```
"..."
```

At the end of a List the"..." string can be placed afterward the last element which means that this last element can be repeated infinitely from here.

This also covers any sub-elements of the previous element, as they are part of that element's spec.

### Indexing Elements

Starting with the first-node elemnent "0", each sub-element adds depth or increments the current counter.

The first child of of "0" will be "0.0".  If there is a peer-element of "0.0", which is the second child of "0", it will be "0.1", third child "0.2", etc.

More children-elements of those continue the depth to "0.1.0" (the second child of the first-node, with a child as well).

The Key Dict type uses the key string instead of the position, so "0.contacts", instead of "0.0".  This is because Key Dicts would require sorting to be positional, which would change if a new key changes the old sorting values.

### Nesting Elements

Using the index numbers in the same way as the Variadic "...", they can be placed in any child-element position (non-first-node), and signifies that the entire structure from that index can appear again in this position.

Placing "0.0.0" would reference that 2nd child object of the first-node.

Example: `"0"`

Example: `{"index": "0"}`

Example: `{"index": "0", "optional": true}`

Three examples showing the same operation, with the final allowing an optional nesting.

### Optional Elements

```
"optional": false
```

The above field is a default value in all of our types.  Any of them can set this to true, and then do not need to exist.  If it is not specified, one of the specified requirements must be met in the data to be valid.

For List elements, only the last element can be optional.  This last element can be variadic, so it can repeat, but having optional elements in the middle of the list is not worth supporting (anti-goal).

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

#### DDD

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

#### DDD

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

The `"*"` character acts as a Glob, so "foo*" and "*foo" or "*foo*" would match their respective strings.  The `"*"` is not allowed inside field names.  If there are conflicts, it chooses the explicit key over a glob, if they are not explicit they are alpha-sorted order, and the first item is chosen.  (NOTE: This could also be the longest key that matches)

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

#### DDD

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

The `"switch_rows"` key can be `"*"` as a glob, and can be used as a global default, or partial glob, in the same ways as the Key Dict globs.

## Just a number

```
6
```

If we want the data to be a single value.

#### DDD

```
{"type": "int", "optional": false}
```

This ensures an integer must exist in this data.  If we make it optional, it could be an empty string or return None/nil/etc, depending on the implementation.

## JSON list of things

```
["Bob", 5.1, [1,2,3,1000]]
```

#### DDD

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

## JSON basic nesting

```
{
    "value": 1234,
    "child":
    {
        "value": 2345,
        "child": {
            "value": 3456
            "child": {}
        }
    }
}
```

This is a single dict `{"value": int, "child": optional_child"}` where the `"optional_child"` is repeating it's parent forever, potentially.  It allows infinite nesting of the root-node's format.

```
{"keydict":
  {
    "value": {"type": "int"}
    "child": {"index": "0", "optional": true}
  }
}
```

The root node is indexed by the "0" value, and so each "child" value can (optionally) have a repeat of the "0" root node in that position.

## JSON non-root-node nesting

```
{
    "contacts": {
        "name": "John Doe",
        "age": 30,
        "phone": "+1-555-555-5555",
        "contacts": {
            "name": "Bob Ross",
            "age": 72,
            "phone": "+1-555-444-7777",
        }
    }
}
```

```
{"keydict":
  {
    "contacts": {"keydict":
      {
        "name": {"type": "string"},
        "age": {"type": "int", "min": 0, "max": 200},
        "phone": {"type": "string"},
        "contacts": {"index": "0.contacts", "optional": true}
      }
    }
  }
}
```

This example is like the previous one, but shows a non-"0" index.  Any position or depth would be the same.


## JSON Advanced Example

```
{
    "message": {
        "_id": 4184,
        "name": "message",
        "fields": {
            "_id": {
                "_id": 17433,
                "name": "_id",
                "field_type": "_int",
                "not_null": true,
                "provision_state": 3
            },
            "data": {
                "_id": 17428,
                "name": "data",
                "field_type": "_document",
                "subfields": {
                    "content": {
                        "_id": 17429,
                        "name": "content",
                        "field_type": "_string",
                        "not_null": true,
                        "provision_state": 3
                    },
                    "created": {
                        "_id": 17432,
                        "name": "created",
                        "field_type": "_int",
                        "not_null": true,
                        "provision_state": 3
                    },
                    "created_by": {
                        "_id": 17431,
                        "name": "created_by",
                        "field_type": "_string",
                        "not_null": true,
                        "provision_state": 3
                    },
                    "thread_id": {
                        "_id": 17430,
                        "name": "thread_id",
                        "field_type": "_int",
                        "not_null": true,
                        "relation": {
                            "_id": 1250,
                            "field_id": 17427,
                            "collection": "thread",
                            "field": "_id"
                        },
                        "provision_state": 3
                    }
                },
                "provision_state": 3
            }
        },
        "indexes": {
            "created": {
                "_id": 5443,
                "name": "created",
                "fields": [
                    "data.created"
                ],
                "provision_state": 3
            }
        },
        "partitions": [{
            "_id": 4180,
            "start_id": 1,
            "shard_config": {
                "shard_key": "_id",
                "hash_method": "cast",
                "shard_method": "mod"
            }
        }],
        "provision_state": 3
    }
}
```

