# Example Terraform file with secure configurations

resource "aws_s3_bucket" "good_bucket" {
  bucket = "my-private-bucket"
}

resource "aws_s3_bucket_public_access_block" "good_bucket" {
  bucket = aws_s3_bucket.good_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_security_group" "good_sg" {
  name_prefix = "secure-security-group"
  
  ingress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    cidr_blocks     = ["10.0.0.0/8"]
    description     = "HTTPS from VPC only"
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }
}