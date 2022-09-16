import asyncio
import sys

from grpc import aio
from grpc import Compression

from federated_train import FederatedTraining
from pythonprotos import federated_pb2 as pb2
from pythonprotos import federated_pb2_grpc as pb2_grpc


class FederatedTrainingService(pb2_grpc.FederatedService):
    def __init__(self):
        self.model_trainer = FederatedTraining()

    def ChangeModeRestart(self, request, context, **kwargs):
        print("Restarting federated server", flush=True)
        sys.exit()

    def TrainUpdateModel(self, request, context, **kwargs):
        try:
            updated_model_weights = self.model_trainer.train_update_model(request)
            updated_model = {"model_id": request.model_id,
                             "layers_weight": updated_model_weights,
                             "client_id": request.client_id,
                             "status": pb2.ModelUpdateRequest.TrainingStatus.SUCCEEDED}
            return pb2.ModelUpdateRequest(**updated_model)
        except Exception as ex:
            print(ex, flush=True)
            updated_model = {"status": pb2.ModelUpdateRequest.TrainingStatus.FAILED}
            return pb2.ModelUpdateRequest(**updated_model)


async def serve_async(port) -> None:
    server = aio.server(options=(
        ("grpc.keepalive_time_ms", 10000),
        ("grpc.keepalive_timeout_ms", 5000),
        ("grpc.keepalive_permit_without_calls", True),
        ("grpc.http2.max_pings_without_data", 0),
        ("grpc.http2.min_time_between_pings_ms", 5000),
        ("grpc.http2.min_ping_interval_without_data_ms", 5000),
        ("grpc.max_connection_age_ms", 600000),
        ("grpc.max_connection_idle_ms", 600000),
        ("grpc.max_connection_age_grace_ms", 30000),
    ),
        compression=Compression.Gzip
    )
    pb2_grpc.add_FederatedServiceServicer_to_server(
        FederatedTrainingService(), server)
    server.add_insecure_port("localhost:" + port)
    await server.start()
    print("Federated server running on port: " + port, flush=True)
    await server.wait_for_termination()


if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(serve_async("9003"))
    # asyncio.run(serve_async("9003"))
