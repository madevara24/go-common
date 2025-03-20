# PostgresMapper

## Overview

PostgresMapper is a custom Go library designed to facilitate dynamic SQL query generation for PostgreSQL databases. It leverages Go's reflection capabilities to automatically adapt to changes in database schema, reducing the need for manual query updates.

## Features

- **Dynamic Query Generation**: Automatically generates SQL queries based on Go struct tags, ensuring alignment with the database schema.
- **CRUD Operations**: Supports basic CRUD operations including select, insert, update, soft-delete, and hard-delete.
- **Batch Processing**: Efficiently handles batch insert and update operations.
- **Soft and Hard Deletes**: Provides flexibility in data management with support for both soft and hard deletes.
- **Field Selection and Filtering**: Allows for precise field selection and query filtering.

## Installation

To install PostgresMapper, use the following command:

```bash
go get github.com/madevara24/go-common/mapper
```

## Usage

Here's a basic example of how to use PostgresMapper:

```go
package main

import (
    "context"
    "github.com/madevara24/go-common/mapper"
)

type User struct {
    ID        int    `db:"id" primarykey:"true"`
    Name      string `db:"name"`
    Email     string `db:"email"`
    CreatedAt time.Time `db:"created_at" hasdbdefault:"true"`
}

func main() {
    mapper := postgresmapper.NewPostgresMapper()
    ctx := context.Background()

    user := User{Name: "John Doe", Email: "john.doe@example.com"}
    query, args, err := mapper.Insert(ctx, user, "users")
    if err != nil {
        panic(err)
    }
    fmt.Println("Generated Query:", query)
    fmt.Println("Arguments:", args)
}
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contact

For questions or support, please contact [adityadevara91@gmail.com].
