package dataconverter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompression(t *testing.T) {
	data := bigJSONPayload()
	compressedData, _ := compressGZV1(data)
	decompressedData, _ := decompressGZV1(compressedData)
	require.Equal(t, len(data), 8591, "original payload byte total is not correct")
	require.True(t, len(compressedData) < 3000, "bytes didn't compress/shrink as expected")
	require.Equal(t, len(data), len(decompressedData), "data did not compress then decompress to original size")
	require.Equal(t, data, decompressedData, "data did not compress then decompress to original data")
}

func bigJSONPayload() []byte {
	return []byte(`[
  {
    "_id": "5ffdd0b0ee4591da4fcef9c3",
    "index": 0,
    "guid": "cd786286-ecf4-4b9e-9b96-b1f4db0df2d0",
    "isActive": false,
    "balance": "$2,825.48",
    "picture": "http://placehold.it/32x32",
    "age": 33,
    "eyeColor": "green",
    "name": "Joseph Mccarty",
    "gender": "male",
    "company": "REPETWIRE",
    "email": "josephmccarty@repetwire.com",
    "phone": "+1 (997) 558-2165",
    "address": "114 Montgomery Street, Goochland, Kentucky, 4208",
    "about": "Aliquip ut officia proident aliqua in deserunt mollit deserunt. Eiusmod laborum ut veniam consequat culpa ipsum ea ut. Ex reprehenderit duis quis irure qui. Nisi sunt id fugiat nisi minim nisi et ad nostrud.\r\n",
    "registered": "2018-03-01T07:08:18 +08:00",
    "latitude": 86.948845,
    "longitude": -92.801158,
    "tags": [
      "ipsum",
      "amet",
      "in",
      "ad",
      "Lorem",
      "aute",
      "culpa"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Cecilia Bryant"
      },
      {
        "id": 1,
        "name": "Mays Barber"
      },
      {
        "id": 2,
        "name": "Peterson Bradshaw"
      }
    ],
    "greeting": "Hello, Joseph Mccarty! You have 4 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "5ffdd0b03463e482de772bae",
    "index": 1,
    "guid": "fa1ccc80-6c10-44de-9580-f8b14cc04fcb",
    "isActive": true,
    "balance": "$1,083.35",
    "picture": "http://placehold.it/32x32",
    "age": 27,
    "eyeColor": "brown",
    "name": "Tammi Cortez",
    "gender": "female",
    "company": "GEEKOLA",
    "email": "tammicortez@geekola.com",
    "phone": "+1 (920) 514-2597",
    "address": "428 Holly Street, Coultervillle, Tennessee, 2741",
    "about": "Proident sint minim laborum commodo enim ut. Et mollit deserunt deserunt eiusmod velit excepteur quis eu adipisicing duis sint nisi proident aute. Dolor mollit voluptate elit mollit ipsum Lorem cillum est dolor aliqua aliqua et eiusmod labore. Irure qui occaecat eu fugiat proident cupidatat cupidatat. Minim esse tempor magna excepteur ullamco sit. Aliquip nulla et qui adipisicing occaecat amet eu tempor laboris anim ipsum excepteur dolor. Tempor fugiat mollit commodo culpa dolor Lorem velit commodo laboris do do laboris.\r\n",
    "registered": "2015-02-20T01:22:30 +08:00",
    "latitude": 52.73284,
    "longitude": 135.911179,
    "tags": [
      "et",
      "anim",
      "consequat",
      "officia",
      "nisi",
      "eu",
      "ea"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Sharlene Kemp"
      },
      {
        "id": 1,
        "name": "Brenda Mcintosh"
      },
      {
        "id": 2,
        "name": "Jerry Perkins"
      }
    ],
    "greeting": "Hello, Tammi Cortez! You have 9 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "5ffdd0b0d3495189d11bd841",
    "index": 2,
    "guid": "9354791c-508b-4fe6-b4fb-ac15f8d5336d",
    "isActive": false,
    "balance": "$2,078.83",
    "picture": "http://placehold.it/32x32",
    "age": 24,
    "eyeColor": "green",
    "name": "Wiggins Rhodes",
    "gender": "male",
    "company": "PURIA",
    "email": "wigginsrhodes@puria.com",
    "phone": "+1 (823) 404-2786",
    "address": "585 John Street, Savage, Wyoming, 9053",
    "about": "Aliquip irure id amet consectetur pariatur labore non irure sit reprehenderit irure enim. Nostrud excepteur id adipisicing adipisicing commodo ea. Esse proident et amet ullamco anim Lorem id commodo. Mollit non magna nisi adipisicing deserunt voluptate laborum occaecat quis excepteur excepteur.\r\n",
    "registered": "2020-01-04T03:00:27 +08:00",
    "latitude": -0.762738,
    "longitude": -25.614977,
    "tags": [
      "nisi",
      "consectetur",
      "sunt",
      "sit",
      "ipsum",
      "amet",
      "mollit"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Gabrielle Chandler"
      },
      {
        "id": 1,
        "name": "Summers Carroll"
      },
      {
        "id": 2,
        "name": "Holmes Pena"
      }
    ],
    "greeting": "Hello, Wiggins Rhodes! You have 3 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "5ffdd0b07d1fc3e0255ed187",
    "index": 3,
    "guid": "c40d4e6b-71ad-4358-95ec-4d61ea84b42f",
    "isActive": true,
    "balance": "$3,130.72",
    "picture": "http://placehold.it/32x32",
    "age": 34,
    "eyeColor": "green",
    "name": "Jenny Grimes",
    "gender": "female",
    "company": "VERBUS",
    "email": "jennygrimes@verbus.com",
    "phone": "+1 (832) 543-2372",
    "address": "306 Malta Street, Allison, Rhode Island, 3694",
    "about": "Consectetur laborum eiusmod eu non excepteur laborum pariatur esse ad proident elit Lorem voluptate. Id est excepteur qui laboris consequat fugiat id sunt non sunt. Ex elit id ex ipsum reprehenderit in anim aliquip aliquip duis laboris duis excepteur officia. Aliquip officia commodo Lorem laboris tempor reprehenderit sunt commodo velit ullamco cupidatat Lorem nulla sit. Ipsum in labore officia ex pariatur occaecat exercitation proident eiusmod nisi fugiat cupidatat incididunt.\r\n",
    "registered": "2019-11-15T05:15:53 +08:00",
    "latitude": 1.774245,
    "longitude": 159.400019,
    "tags": [
      "ad",
      "ea",
      "enim",
      "excepteur",
      "nostrud",
      "velit",
      "nisi"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Ellison Stein"
      },
      {
        "id": 1,
        "name": "Emilia Oconnor"
      },
      {
        "id": 2,
        "name": "Luann Wright"
      }
    ],
    "greeting": "Hello, Jenny Grimes! You have 10 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "5ffdd0b047f1d0e30187949a",
    "index": 4,
    "guid": "bdfa50e3-94d2-4d86-96ce-df7d50b65887",
    "isActive": false,
    "balance": "$2,991.15",
    "picture": "http://placehold.it/32x32",
    "age": 30,
    "eyeColor": "green",
    "name": "Bush Tucker",
    "gender": "male",
    "company": "ANIMALIA",
    "email": "bushtucker@animalia.com",
    "phone": "+1 (808) 572-3654",
    "address": "332 Fanchon Place, Hampstead, Hawaii, 5437",
    "about": "Cupidatat nulla ex duis sit consectetur eiusmod mollit eiusmod labore. Veniam ipsum quis anim excepteur do esse. Ut ad culpa fugiat velit esse ad eiusmod laboris culpa magna eiusmod. Tempor anim cupidatat ipsum aute ex excepteur non commodo magna cupidatat deserunt aliqua nulla amet. Mollit pariatur adipisicing dolore dolor labore ipsum aliquip deserunt.\r\n",
    "registered": "2015-04-02T06:25:30 +07:00",
    "latitude": -6.753237,
    "longitude": 48.24015,
    "tags": [
      "amet",
      "consequat",
      "esse",
      "sunt",
      "velit",
      "reprehenderit",
      "dolore"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Sandy Young"
      },
      {
        "id": 1,
        "name": "Guzman Kramer"
      },
      {
        "id": 2,
        "name": "Monroe Cash"
      }
    ],
    "greeting": "Hello, Bush Tucker! You have 6 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "5ffdd0b00e56c831a7b3584e",
    "index": 5,
    "guid": "779f1114-98d5-48e6-b43c-d3bb6faeedf2",
    "isActive": true,
    "balance": "$2,141.15",
    "picture": "http://placehold.it/32x32",
    "age": 20,
    "eyeColor": "brown",
    "name": "Jessica Foster",
    "gender": "female",
    "company": "VERTIDE",
    "email": "jessicafoster@vertide.com",
    "phone": "+1 (957) 454-2713",
    "address": "945 Royce Street, Fairacres, New Hampshire, 6405",
    "about": "Mollit minim ea cupidatat ullamco cillum adipisicing aliqua minim ex labore. Et culpa nulla voluptate culpa cupidatat enim officia velit amet officia excepteur. Aliqua sint ea adipisicing pariatur est amet. Quis excepteur occaecat est fugiat ex occaecat do anim et nisi voluptate dolore. Occaecat mollit consequat culpa reprehenderit minim nostrud minim sit et eu. Consequat esse proident proident elit quis officia nulla et. Aliquip fugiat sunt reprehenderit proident voluptate.\r\n",
    "registered": "2020-10-07T09:56:16 +07:00",
    "latitude": 81.833425,
    "longitude": 107.252947,
    "tags": [
      "tempor",
      "aliquip",
      "ex",
      "adipisicing",
      "nostrud",
      "consectetur",
      "ex"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Donaldson Peterson"
      },
      {
        "id": 1,
        "name": "Estella Morton"
      },
      {
        "id": 2,
        "name": "Downs Cherry"
      }
    ],
    "greeting": "Hello, Jessica Foster! You have 7 unread messages.",
    "favoriteFruit": "banana"
  }
]`)
}

