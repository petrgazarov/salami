resource "aws_vpc" "MainVpc" {
  cidr_block       = "10.0.0.0/16"
  enable_dns_support = true
  enable_dns_hostnames = true

  tags = {
    Name = "main-vpc"
  }
}

resource "aws_subnet" "PrivateSubnetA" {
  vpc_id                  = aws_vpc.MainVpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-west-1a"
  map_public_ip_on_launch = false
  tags = {
    Name = "private-subnet-a"
  }
}

resource "aws_subnet" "PrivateSubnetB" {
  vpc_id                     = aws_vpc.MainVpc.id
  cidr_block                 = "10.0.2.0/24"
  availability_zone          = "us-west-1c"
  map_public_ip_on_launch    = false
  tags = {
    Name = "private-subnet-b"
  }
}

resource "aws_route_table" "PrivateRouteTable" {
  vpc_id = aws_vpc.MainVpc.id

  tags = {
    Name = "private-route-table"
  }
}

resource "aws_route_table_association" "PrivateSubnetARouteTableAssociation" {
  subnet_id      = aws_subnet.PrivateSubnetA.id
  route_table_id = aws_route_table.PrivateRouteTable.id
}

resource "aws_route_table_association" "PrivateSubnetBRouteTableAssociation" {
  subnet_id      = aws_subnet.PrivateSubnetB.id
  route_table_id = aws_route_table.PrivateRouteTable.id
}

resource "aws_subnet" "PublicSubnetA" {
  vpc_id                     = aws_vpc.MainVpc.id
  cidr_block                 = "10.0.3.0/24"
  availability_zone          = "us-west-1a"
  map_public_ip_on_launch    = true

  tags = {
    Name = "public-subnet-a"
  }
}

resource "aws_subnet" "PublicSubnetB" {
  vpc_id                     = aws_vpc.MainVpc.id
  cidr_block                 = "10.0.4.0/24"
  availability_zone          = "us-west-1c"
  map_public_ip_on_launch    = true

  tags = {
    Name = "public-subnet-b"
  }
}

resource "aws_internet_gateway" "InternetGateway" {
  vpc_id = aws_vpc.MainVpc.id

  tags = {
    Name = "internet-gateway"
  }
}

resource "aws_route_table" "PublicRouteTable" {
  vpc_id = aws_vpc.MainVpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.InternetGateway.id
  }

  tags = {
    Name = "public-route-table"
  }
}

resource "aws_route_table_association" "PublicSubnetARouteTableAssociation" {
  subnet_id      = aws_subnet.PublicSubnetA.id
  route_table_id = aws_route_table.PublicRouteTable.id
}

resource "aws_route_table_association" "PublicSubnetBRouteTableAssociation" {
  subnet_id      = aws_subnet.PublicSubnetB.id
  route_table_id = aws_route_table.PublicRouteTable.id
}