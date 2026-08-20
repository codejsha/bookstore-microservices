terraform {
  required_providers {
    opensearch = {
      source = "opensearch-project/opensearch"
    }
  }
}

resource "opensearch_index_template" "books" {
  name = "books"
  body = jsonencode({
    index_patterns = ["books", "books-*"]
    priority       = 100
    template = {
      settings = {
        number_of_shards         = 3
        number_of_replicas       = 1
        "index.refresh_interval" = "2s"
      }
      mappings = {
        dynamic = "strict"
        properties = {
          uid = { type = "keyword" }
          title = { type = "text", analyzer = "standard",
          fields = { raw = { type = "keyword" } } }
          description        = { type = "text", analyzer = "standard" }
          ol_key             = { type = "keyword" }
          first_publish_date = { type = "keyword" }
          cover_uids         = { type = "keyword" }
          created_at         = { type = "date" }
          updated_at         = { type = "date" }

          authors = {
            type = "nested"
            properties = {
              uid             = { type = "keyword" }
              name            = { type = "text", fields = { raw = { type = "keyword" } } }
              bio             = { type = "text" }
              birth_date      = { type = "keyword" }
              death_date      = { type = "keyword" }
              ol_key          = { type = "keyword" }
              alternate_names = { type = "keyword" }
            }
          }

          subjects = {
            type = "nested"
            properties = {
              uid  = { type = "keyword" }
              name = { type = "text", fields = { raw = { type = "keyword" } } }
            }
          }

          editions = {
            type = "nested"
            properties = {
              uid             = { type = "keyword" }
              title           = { type = "text", fields = { raw = { type = "keyword" } } }
              isbn10          = { type = "keyword" }
              isbn13          = { type = "keyword" }
              publish_date    = { type = "keyword" }
              languages       = { type = "keyword" }
              physical_format = { type = "keyword" }
              description     = { type = "text" }
              publisher = {
                type = "object"
                properties = {
                  uid  = { type = "keyword" }
                  name = { type = "text", fields = { raw = { type = "keyword" } } }
                }
              }
            }
          }
        }
      }
    }
  })
}

resource "opensearch_index" "books_v1" {
  depends_on = [opensearch_index_template.books]
  name       = "books-v1"
  aliases = jsonencode({
    books = {
      is_write_index = true
    }
  })
}
