#!/bin/bash

protoc grpc-calculator/proto/calc.proto --go_out=plugins=grpc:.

protoc greet/proto/greet.proto --go_out=plugins=grpc:.