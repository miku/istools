istools
=======

* islint, an intermediate schema linter (record quality checks)
* islabel, a sigel attacher (determine license coverage of records)

Sketches
--------

Licensing tree.

```json
{
  "A": {
    "match_all": {}
  },
  "B": {
    "holding": {
      "location": "/path/to/file"
    }
  },
  "C": {
    "or": [
      {
        "holding": {
          "location": "/path/to/file"
        }
      },
      {
        "attr": {
          "path": "a.b.c",
          "regex": ".*XYZ.*"
        }
      }
    ]
  },
  "D": {
    "or": [
      {
        "and": [
          {
            "attr": {
              "path": "x.a",
              "list": "/path/to/file"
            }
          },
          {
            "attr": {
              "path": "x.b",
              "value": "123"
            }
          }
        ]
      },
      {
        "holding": {
          "location": "/path/to/file"
        }
      },
      {
        "attr": {
          "path": "a.b.c",
          "regex": ".*XYZ.*"
        }
      }
    ]
  }
}
```