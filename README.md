# R2 Action

Perform an operation with objects in Cloudflare R2 bucket

## Inputs

| Name              | Description                                        | Required |
|-------------------|----------------------------------------------------|----------|
| account_id        | Cloudflare Account ID                              | true     |
| access_key_id     | Cloudflare R2 Access Key ID                        | true     |
| access_key_secret | Cloudflare R2 Access Key Secret                    | true     |
| bucket            | Cloudflare R2 Bucket Name                          | true     |
| operation         | Operation to perform: get, put, delete             | true     |
| key               | Object key                                         | true     |
| file              | Local file path, not required for delete operation | false    |
