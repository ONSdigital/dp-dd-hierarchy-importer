package geography

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var OriginalJson = `{
  "ons": {
    "geographyList": {
      "geography": {
        "id": "2011STATH",
        "names": {
          "name": [
            {
              "@xml.lang": "en",
              "$": "2011 Statistical Geography Hierarchy"
            },
            {
              "@xml.lang": "cy",
              "$": "Hierarchaeth Daearyddiaeth Ystadegol 2011"
            }
          ]
        }
      },
      "items": {
        "item": [
          {
            "labels": {
              "label": [
                {
                  "@xml.lang": "en",
                  "$": "United Kingdom"
                },
                {
                  "@xml.lang": "cy",
                  "$": "Deyrnas Unedig"
                }
              ]
            },
            "itemCode": "K02000001",
            "areaType": {
              "abbreviation": "UK",
              "codename": "United Kingdom",
              "level": 0
            },
            "subthresholdAreas": ""
          },
          {
            "labels": {
              "label": [
                {
                  "@xml.lang": "en",
                  "$": "England and Wales"
                },
                {
                  "@xml.lang": "cy",
                  "$": "Cymru a Lloegr"
                }
              ]
            },
            "itemCode": "K04000001",
            "parentCode": "K03000001",
            "areaType": {
              "abbreviation": "NAT",
              "codename": "England and Wales",
              "level": 2
            },
            "subthresholdAreas": ""
          },
          {
            "labels": {
              "label": [
                {
                  "@xml.lang": "en",
                  "$": "Great Britain"
                },
                {
                  "@xml.lang": "cy",
                  "$": "Prydain Fawr"
                }
              ]
            },
            "itemCode": "K03000001",
            "parentCode": "K02000001",
            "areaType": {
              "abbreviation": "GB",
              "codename": "Great Britain",
              "level": 1
            },
            "subthresholdAreas": ""
          }
        ]
      }
    }
  }
}`

func TestMarshalAndUnmarshal(t *testing.T) {
	Convey("When json is unmarshaled", t, func() {
		var data *GeographicData
		err := json.Unmarshal([]byte(OriginalJson), &data)

		So(err, ShouldBeNil)

		Convey("When json is marshaled again", func() {
			newJson, err := json.MarshalIndent(&data, "", "  ")

			So(err, ShouldBeNil)
			So(string(newJson), ShouldEqual, OriginalJson)
		})
	})
}
