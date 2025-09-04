db = db.getSiblingDB('movies_db');

// Create movies collection with schema validation
db.createCollection('movies', {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["_id", "title", "year"],
         properties: {
            _id: {
               bsonType: "int",
               description: "must be an integer and is required"
            },
            title: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            year: {
               bsonType: "string",
               pattern: "^[0-9]{4}$",
               description: "must be a 4-digit year string and is required"
            }
         }
      }
   }
});

// Create index for text search and year
db.movies.createIndex({ "title": "text", "year": 1 });

print("MongoDB initialization completed successfully!");