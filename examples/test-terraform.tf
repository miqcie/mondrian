# Example Terraform file with security violations for testing

resource "aws_s3_bucket" "bad_bucket" {
  bucket = "my-public-bucket"
}

resource "aws_s3_bucket_public_access_block" "bad_bucket" {
  bucket = aws_s3_bucket.bad_bucket.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "bad_bucket_acl" {
  bucket = aws_s3_bucket.bad_bucket.id
  acl    = "public-read"
}

resource "aws_security_group" "bad_sg" {
  name_prefix = "bad-security-group"
  
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}