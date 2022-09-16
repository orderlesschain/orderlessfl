import os

import numpy
import tensorflow as tf

from pythonprotos import federated_pb2 as pb2


class FederatedTraining:
    def __init__(self):
        os.environ["TF_CPP_MIN_LOG_LEVEL"] = "2"
        (train_images, train_labels), (test_images, test_labels) = tf.keras.datasets.mnist.load_data()
        self.train_labels = train_labels
        self.test_labels = test_labels
        self.train_images = train_images / 255.0
        self.test_images = test_images / 255.0

    def train_update_model(self, update_request):
        model = FederatedTraining.create_model(update_request.model_type)
        server_weights = []
        model_clock = 0
        if update_request.layers_weight is not None and len(update_request.layers_weight) > 0:
            (weights, model_clock) = FederatedTraining.convert_proto_to_model_weights(update_request.layers_weight)
            server_weights.extend(weights)
        else:
            print("no model weight from server", flush=True)
            server_weights.extend(FederatedTraining.read_init_model(update_request.model_id))
        model.set_weights(server_weights)
        ss, acc = model.evaluate(self.test_images, self.test_labels, verbose=0)
        print("accuracy BEFORE training: {:5.2f}%".format(100 * acc), flush=True)
        beginning, end = self.get_range_for_training_data(update_request.total_clients_count,
                                                          update_request.batch_size,
                                                          update_request.client_id,
                                                          update_request.request_counter)
        image_length = len(self.train_images)
        if end > image_length:
            end = image_length
            beginning = image_length - update_request.batch_size
        model.fit(self.train_images[beginning:end],
                  self.train_labels[beginning:end],
                  batch_size=128,
                  verbose=0,
                  epochs=10)
        ss, acc = model.evaluate(self.test_images, self.test_labels, verbose=0)
        print("accuracy AFTER training: {:5.2f}%".format(100 * acc), flush=True)
        trained_model_diff = FederatedTraining.create_delta_weights(model.get_weights(), server_weights)
        # print("Model Clock " + str(model_clock), flush=True)

        # model3 = FederatedTraining.create_model()
        # # ss, acc = model3.evaluate(self.test_images, self.test_labels, verbose=0)
        # # print("Model, accuracy: {:5.2f}%".format(100 * acc), flush=True)
        #
        # # new_weight = FederatedTraining.add_delta_weights(server_weights, trained_model_diff)
        # new_weight = FederatedTraining.add_delta_weights(trained_model_diff, server_weights)
        #
        # model3.set_weights(new_weight)
        #
        # ss, acc = model3.evaluate(self.test_images, self.test_labels, verbose=0)
        # print("Model, accuracy: {:5.2f}%".format(100 * acc), flush=True)

        return FederatedTraining.convert_model_weights_to_proto(trained_model_diff, model_clock,
                                                                update_request.client_id,
                                                                update_request.client_latest_clock)

    def get_range_for_training_data(self, total_clients_count, batch_size, client_id, request_counter):
        client_range = len(self.train_images) / total_clients_count
        beginning = (client_id * client_range) + (batch_size * request_counter)
        end = beginning + batch_size
        return int(beginning), int(end)

    @staticmethod
    def create_delta_weights(locally_trained_model, server_model):
        if len(locally_trained_model) != len(server_model):
            raise Exception("the shapes of models are not matching")
        delta_weights = []
        for index in range(len(server_model)):
            delta_weights.append(numpy.subtract(locally_trained_model[index], server_model[index]))
        return delta_weights

    @staticmethod
    def add_delta_weights(old_model, delta):
        if len(old_model) != len(delta):
            raise Exception("the shapes of models are not matching")
        delta_weights = []
        for index in range(len(old_model)):
            delta_weights.append(numpy.add(old_model[index], delta[index]))
        return delta_weights

    @staticmethod
    def create_model(model_type_given):
        model_parameter_count = 0
        if model_type_given == 1:
            model_parameter_count = 128
        elif model_type_given == 2:
            model_parameter_count = 256
        elif model_type_given == 3:
            model_parameter_count = 384
        elif model_type_given == 4:
            model_parameter_count = 512
        model = tf.keras.models.Sequential([
            tf.keras.layers.Flatten(input_shape=(28, 28)),
            tf.keras.layers.Dense(model_parameter_count, activation='relu'),
            tf.keras.layers.Dropout(0.2),
            tf.keras.layers.Dense(10)
        ])
        model.compile(optimizer='adam',
                      loss=tf.keras.losses.SparseCategoricalCrossentropy(from_logits=True),
                      metrics=['accuracy'])

        return model

    @staticmethod
    def convert_model_weights_to_proto(updated_weights, model_clock, client_id, client_latest_clock):
        proto_converted_model = pb2.ModelWeights()
        proto_converted_model.model_clock = model_clock
        proto_converted_model.client_id = client_id
        proto_converted_model.client_latest_clock = client_latest_clock
        index = 0
        for weight in updated_weights:
            layer_weights = pb2.LayerWeights()
            layer_weights.layer_rank = index
            if len(weight.shape) == 1:
                layer_weights.weight_type = pb2.LayerWeights.LayerWeightsType.ONE_D
                layer_weights.one_d_weights.extend(weight.tolist())
            elif len(weight.shape) == 2:
                layer_weights.weight_type = pb2.LayerWeights.LayerWeightsType.TWO_D
                weight_2_ds = weight.tolist()
                for weight_2_d in weight_2_ds:
                    proto_weights = pb2.LayerTwoDWeights()
                    proto_weights.weights.extend(weight_2_d)
                    layer_weights.two_d_weights.append(proto_weights.SerializeToString())
            else:
                raise Exception("weight to proto, because the weight shape is ", len(weight.shape))
            proto_converted_model.layers.append(layer_weights.SerializeToString())
            index += 1
        string_weights = proto_converted_model.SerializeToString()
        # print("Client Id " + str(proto_converted_model.client_id), flush=True)

        return string_weights

    @staticmethod
    def convert_proto_to_model_weights(string_weights):
        model_weights_dict = {}
        proto_converted_model = pb2.ModelWeights()
        proto_converted_model.ParseFromString(string_weights)
        model_clock = proto_converted_model.model_clock
        for layer_weights_proto in proto_converted_model.layers:
            layer_weights = pb2.LayerWeights()
            layer_weights.ParseFromString(layer_weights_proto)
            if layer_weights.weight_type == pb2.LayerWeights.LayerWeightsType.ONE_D:
                model_weights_dict[layer_weights.layer_rank] = []
                model_weights_dict[layer_weights.layer_rank].extend(numpy.asarray(layer_weights.one_d_weights))
            elif layer_weights.weight_type == pb2.LayerWeights.LayerWeightsType.TWO_D:
                temp_2_d_weights = []
                model_weights_dict[layer_weights.layer_rank] = []
                for weight_2_ds in layer_weights.two_d_weights:
                    layer_2_proto_weights = pb2.LayerTwoDWeights()
                    layer_2_proto_weights.ParseFromString(weight_2_ds)
                    temp_2_d_weights.append(numpy.asarray(layer_2_proto_weights.weights))
                model_weights_dict[layer_weights.layer_rank].extend(numpy.asarray(temp_2_d_weights))
            else:
                raise Exception("proto to weight error, because the layer type is ",
                                layer_weights.weight_type)
        model_weights = []
        for index in range(len(model_weights_dict)):
            model_weights.append(numpy.asarray(model_weights_dict[index]))
        return model_weights, model_clock

    @staticmethod
    def read_init_model(model_id):
        server_weights = []
        try:
            f = open("./federated_models/" + model_id, "rb")
            (weights, model_clock) = FederatedTraining.convert_proto_to_model_weights(f.read())
            server_weights.extend(weights)
            f.close()
        except IOError:
            print("could not open init model", flush=True)
        return server_weights


if __name__ == "__main__":
    pass
