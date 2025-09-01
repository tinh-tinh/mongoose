# Mongoose for Tinh Tinh

<div align="center">
<img alt="GitHub Release" src="https://img.shields.io/github/v/release/tinh-tinh/mongoose">
<img alt="GitHub License" src="https://img.shields.io/github/license/tinh-tinh/mongoose">
<a href="https://codecov.io/gh/tinh-tinh/mongoose" > 
 <img src="https://codecov.io/gh/tinh-tinh/mongoose/branch/master/graph/badge.svg?token=EP4XOF5HOY"/> 
</a>
<a href="https://pkg.go.dev/github.com/tinh-tinh/mongoose"><img src="https://pkg.go.dev/badge/github.com/tinh-tinh/mongoose.svg" alt="Go Reference"></a>
</div>

<div align="center">
    <img src="https://avatars.githubusercontent.com/u/178628733?s=400&u=2a8230486a43595a03a6f9f204e54a0046ce0cc4&v=4" width="200" alt="Tinh Tinh Logo">
</div>

## Overview

Mongoose for Tinh Tinh is a powerful MongoDB integration package designed specifically for the Tinh Tinh framework. It provides a clean, efficient, and type-safe way to work with MongoDB databases in your Go applications.

## Features

- 🚀 Simple and intuitive MongoDB integration
- 📦 BSON document handling
- 🔄 Automatic model serialization/deserialization
- 🎯 Type-safe query building
- 🛠️ Advanced features including:
  - Collection operations
  - Aggregation pipelines
  - Index management
  - Transaction support
  - Change streams
  - GridFS support

## Installation

To install the package, use:

```bash
go get -u github.com/tinh-tinh/mongoose/v2
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/tinh-tinh/mongoose/v2"
)

// User represents your MongoDB document
type User struct {
    ID       string `bson:"_id,omitempty"`
    Name     string `bson:"name"`
    Email    string `bson:"email"`
    Age      int    `bson:"age"`
    IsActive bool   `bson:"is_active"`
}

func main() {
    // Initialize MongoDB connection
    client := mongoose.New(&mongoose.Config{
        URI:      "mongodb://localhost:27017",
        Database: "myapp",
    })

    // Get a collection
    collection := client.Collection("users")

    // Insert a document
    user := User{
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      30,
        IsActive: true,
    }
    
    result, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        panic(err)
    }
}
```

## Configuration

The package supports various configuration options:

```go
type Config struct {
    URI              string        // MongoDB connection URI
    Database         string        // Database name
    MaxPoolSize      uint64        // Maximum number of connections
    MinPoolSize      uint64        // Minimum number of connections
    ConnectTimeout   time.Duration // Connection timeout
    MaxConnIdleTime  time.Duration // Maximum idle connection time
    RetryWrites     bool          // Enable retry writes
    RetryReads      bool          // Enable retry reads
    DirectConnection bool          // Use direct connection
}
```

## Key Features

### Collection Operations
```go
// Find documents
users, err := collection.Find(ctx, bson.M{"age": bson.M{"$gt": 25}})

// Update documents
update := bson.M{"$set": bson.M{"is_active": false}}
result, err := collection.UpdateMany(ctx, bson.M{"age": bson.M{"$lt": 18}}, update)

// Delete documents
result, err := collection.DeleteOne(ctx, bson.M{"email": "john@example.com"})
```

### Aggregation Pipeline
```go
pipeline := mongo.Pipeline{
    {{$match: {"age": {"$gt": 25}}}},
    {{$group: {"_id": "$city", "count": {"$sum": 1}}}}
}
cursor, err := collection.Aggregate(ctx, pipeline)
```

### Transactions
```go
err := client.UseTransaction(ctx, func(sessCtx mongo.SessionContext) error {
    // Perform operations within transaction
    _, err := collection.InsertOne(sessCtx, newUser)
    if err != nil {
        return err
    }
    _, err = collection.UpdateOne(sessCtx, filter, update)
    return err
})
```

## Best Practices

1. **Connection Management**
   - Always close connections when done
   - Use appropriate pool sizes
   - Handle connection errors properly

2. **Error Handling**
   - Check for both operation and connection errors
   - Implement proper retry logic
   - Log relevant error information

3. **Performance Optimization**
   - Use appropriate indexes
   - Implement efficient queries
   - Monitor query performance

## Documentation

For detailed documentation and examples, please visit:
- [Go Package Documentation](https://pkg.go.dev/github.com/tinh-tinh/mongoose)
- [MongoDB Go Driver Documentation](https://docs.mongodb.com/drivers/go)

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

If you encounter any issues or need help, you can:
- Open an issue in the GitHub repository
- Check our documentation
- Join our community discussions
