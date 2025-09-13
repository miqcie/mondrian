# This file intentionally contains security violations for testing

resource "aws_s3_bucket" "public_bucket" {
  bucket = "my-public-test-bucket"
}

resource "aws_s3_bucket_acl" "public_bucket_acl" {
  bucket = aws_s3_bucket.public_bucket.id
  acl    = "public-read"
}

resource "aws_security_group" "open_sg" {
  name_prefix = "open-security-group"
  
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH from anywhere - SECURITY VIOLATION!"
  }
}