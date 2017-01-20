package geography

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// OrignalJSON extracted from: http://web.ons.gov.uk/ons/api/data/hierarchies/hierarchy/2011STATH.json?apikey=API_KEY&levels=1,2,3,4
var OriginalJSON = `{
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
          },
          {
            "labels": {
              "label": [
                {
                  "@xml.lang": "en",
                  "$": "England"
                },
                {
                  "@xml.lang": "cy",
                  "$": "Lloegr"
                }
              ]
            },
            "itemCode": "E92000001",
            "parentCode": "K04000001",
            "areaType": {
              "abbreviation": "CTRY",
              "codename": "Country",
              "level": 3
            },
            "subthresholdAreas": ""
          },
          {
            "labels": {
              "label": [
                {
                  "@xml.lang": "en",
                  "$": "Wales"
                },
                {
                  "@xml.lang": "cy",
                  "$": "Cymru"
                }
              ]
            },
            "itemCode": "W92000004",
            "parentCode": "K04000001",
            "areaType": {
              "abbreviation": "CTRY",
              "codename": "Country",
              "level": 3
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
		err := json.Unmarshal([]byte(OriginalJSON), &data)

		So(err, ShouldBeNil)

		Convey("When json is marshaled again", func() {
			newJSON, err := json.MarshalIndent(&data, "", "  ")

			So(err, ShouldBeNil)
			So(string(newJSON), ShouldEqual, OriginalJSON)
		})
	})
}
