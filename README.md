# xgor

`xgor` is a library that extends [Gorm](https://gorm.io/) to provide additional functionalities for building robust database repositories with support for custom filters, transactions, and relationship handling.

## Features

- **Generic Repository:** Use generic repository patterns to handle common CRUD operations for your Gorm models.

- **Custom Filters:** Easily filter entities based on custom conditions using a flexible and intuitive filter syntax.

- **Transaction Support:** Perform operations within a transaction to ensure consistency and atomicity.

- **Relationship Handling:** Simplify relationship management with built-in functions for clearing relationships.

## Installation

```bash
go get -u github.com/grahms/xgor
```

## Usage

### Initializing a Repository

```go
import (
    "gorm.io/gorm"
    "github.com/grahms/xgor"
)

// Initialize  DB
db, err := xgor.Open(...)

// Create a new repository
repo := xgor.New[BlogPost](db, errors.New("blog post not found"))

// Or with relationships
repoWithRelations := xgor.NewWithRelationships[BlogPost](db, errors.New("blog post not found"), "comments", "author")
```

### Adding a Blog Post

```go
post := &BlogPost{
    Title:       "Introduction to xgor",
    Content:     "Learn how to use xgor to supercharge your Gorm-based repositories.",
    AuthorID:    1,
    CategoryID:   2,
    PublishedAt:  time.Now(),
}
err := repo.Add(post)
```

### Querying Blog Posts with Custom Filters

```go
// Get all published posts in the "Technology" category written by a specific author
filters := xgor.FilterType{
    "published_at__lte": time.Now(),
    "category.name__eq": "Technology",
    "author.id__eq":     1,
}
posts, err := repo.GetAll(nil, nil, nil, filters)
```

### Performing a Transaction

```go
err := repo.PerformTransaction(func(tx *gorm.DB) error {
    // Update the author's profile and add a new blog post within the same transaction
    author, err := authorRepo.GetByID(1)
    if err != nil {
        return err
    }

    author.Name = "Updated Author Name"
    if err := authorRepo.Update(author); err != nil {
        return err
    }

    newPost := &BlogPost{
        Title:       "Advanced xgor Techniques",
        Content:     "Explore advanced techniques for optimizing database queries with xgor.",
        AuthorID:    1,
        CategoryID:   3,
        PublishedAt:  time.Now(),
    }

    return repo.Add(newPost)
})
```

## Example Use Case: Blogging Application

Let's consider a blogging application where `xgor` is used to manage blog posts. In this scenario, `xgor` simplifies the data access layer, allowing developers to focus on building features rather than dealing with intricate database operations.

### Use Case Scenario

- **Scenario:** The application needs to fetch all published blog posts in a specific category written by a particular author.

- **Solution:** Utilize `xgor`'s custom filters to easily query the database and retrieve the required blog posts without the complexity of crafting intricate SQL queries.

```go
// Example: Get all published posts in the "Technology" category written by a specific author
filters := xgor.FilterType{
    "published_at__lte": time.Now(),
    "category.name__eq": "Technology",
    "author.id__eq":     1,
}
posts, err := repo.GetAll(nil, nil, nil, filters)
```
## Example Use Case: Blogging Application (Pagination)

### Use Case Scenario

- **Scenario:** The blogging application needs to display a paginated list of blog posts on the homepage.

- **Solution:** Utilize `xgor` to implement pagination and retrieve a subset of blog posts for display.

```go
// Example: Get paginated blog posts for the homepage
limit := 10  // Number of posts per page
page := 1    // Current page
orderBy := "published_at desc"  // Order posts by published date in descending order

// Use xgor to get paginated blog posts
paginationFilters := xgor.FilterType{"category_id__eq": 1}  // Filter by category ID, if needed
blogPosts, err := repo.GetAll(&limit, &page, &orderBy, paginationFilters)

// Check for errors and handle the paginated blog posts
if err != nil {
    // Handle error
} else {
    // Access paginated results
    totalPosts := blogPosts.TotalCount
    currentPage := page
    postsPerPage := limit
    resultCount := blogPosts.ResultCount
    displayedPosts := *blogPosts.Items

    // Process and display paginated blog posts
    for _, post := range displayedPosts {
        // Process each blog post
    }
}
```
## Custom Filters

Custom filters allow you to specify conditions for filtering entities. The filter syntax is based on the column name and a suffix that represents the condition. Here is a table of available filters:

| Filter          | Description                                  | Example                           |
|------------------|----------------------------------------------|-----------------------------------|
| `__eq`           | Equals                                       | `"age__eq": 25`                   |
| `__gt`           | Greater Than                                  | `"age__gt": 21`                   |
| `__lt`           | Less Than                                     | `"age__lt": 30`                   |
| `__gte`          | Greater Than or Equal To                      | `"age__gte": 21`                  |
| `__lte`          | Less Than or Equal To                         | `"age__lte": 30`                  |
| `__in`           | In Array                                      | `"age__in": []int{25, 30}`        |
| `__not`          | Not Equal To                                  | `"age__not": 25`                  |
| `__not_in`       | Not In Array                                  | `"age__not_in": []int{25, 30}`    |
| `__like`         | Like (substring match)                        | `"name__like": "John"`            |

Combine these filters to create powerful and flexible queries tailored to your application's needs.

## Contributing

Feel free to contribute by opening issues or submitting pull requests. Please follow the [Contributing Guidelines](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
