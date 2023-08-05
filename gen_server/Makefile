


pb/chat_pb2_grpc.py pb/chat_pb2.py pb/chat_pb2.pyi:
	@pip install grpcio-tools==1.51.1 mypy-protobuf==3.4.0 'types-protobuf>=3.20.4' --no-cache-dir
	@mkdir -p  ./pb || true
	@echo "===更新生成 chat.proto python 端代码==="
	@python -m grpc_tools.protoc  -I../llms/local/  --python_out=./pb  --pyi_out=./pb  --grpc_python_out=./pb ../llms/local/chat.proto


proto: pb/chat_pb2_grpc.py pb/chat_pb2.py pb/chat_pb2.pyi

