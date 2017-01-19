package structure

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var OriginalJson = `{
  "Structure": {
    "CodeLists": {
      "CodeList": [
        {
          "@id": "CL_0000641",
          "Name": [
            {
              "@xml.lang": "en",
              "$": "COICOP"
            },
            {
              "@xml.lang": "cy"
            }
          ],
          "Code": [
            {
              "@value": "CI_0004280",
              "@urn": "",
              "@parentCode": "CI_0004227",
              "Description": {
                "@xml.lang": "en",
                "$": "05.1.1 Furniture and furnishings"
              },
              "Annotations": {
                "Annotation": [
                  {
                    "AnnotationType": "SubTotal",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "0"
                    }
                  },
                  {
                    "AnnotationType": "IsTotal",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "0"
                    }
                  },
                  {
                    "AnnotationType": "DisplayOrder",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "1"
                    }
                  }
                ]
              }
            },
            {
              "@value": "CI_0005201",
              "@urn": "",
              "@parentCode": "",
              "Description": {
                "@xml.lang": "en",
                "$": "942000 CPI excluding Energy, Food, Alcohol & Tobacco"
              },
              "Annotations": {
                "Annotation": [
                  {
                    "AnnotationType": "SubTotal",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "0"
                    }
                  },
                  {
                    "AnnotationType": "IsTotal",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "0"
                    }
                  },
                  {
                    "AnnotationType": "DisplayOrder",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "75"
                    }
                  }
                ]
              }
            }
          ]
        },
        {
          "@id": "CL_0000641",
          "Name": [
            {
              "@xml.lang": "en",
              "$": "Special Aggregate"
            },
            {
              "@xml.lang": "cy"
            }
          ],
          "Code": [
            {
              "@value": "CI_0004267",
              "@urn": "",
              "@parentCode": "CI_0004221",
              "Description": {
                "@xml.lang": "en",
                "$": "03.1.2 Garments"
              },
              "Annotations": {
                "Annotation": [
                  {
                    "AnnotationType": "SubTotal",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "0"
                    }
                  },
                  {
                    "AnnotationType": "IsTotal",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "0"
                    }
                  },
                  {
                    "AnnotationType": "DisplayOrder",
                    "AnnotationText": {
                      "@xml.lang": "en",
                      "$": "3"
                    }
                  }
                ]
              }
            }
          ]
        }
      ]
    }
  }
}`

func TestMarshalAndUnmarshal(t *testing.T) {
	Convey("When json is unmarshaled", t, func() {
		var data *StructuralData
		err := json.Unmarshal([]byte(OriginalJson), &data)

		So(err, ShouldBeNil)

		Convey("When json is marshaled again", func() {
			newJson, err := json.MarshalIndent(&data, "", "  ")
			newJson = bytes.Replace(newJson, []byte("\\u0026"), []byte("&"), -1)

			So(err, ShouldBeNil)
			So(string(newJson), ShouldEqual, OriginalJson)
		})
	})
}
