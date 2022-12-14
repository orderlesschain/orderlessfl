# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: federated.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import shared_pb2 as shared__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='federated.proto',
  package='protos',
  syntax='proto3',
  serialized_options=b'Z\t./;protos',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0f\x66\x65\x64\x65rated.proto\x12\x06protos\x1a\x0cshared.proto\"\xb3\x02\n\x12ModelUpdateRequest\x12\x10\n\x08model_id\x18\x01 \x01(\t\x12\x12\n\nmodel_type\x18\n \x01(\x05\x12\x15\n\rlayers_weight\x18\x02 \x01(\x0c\x12\x39\n\x06status\x18\x03 \x01(\x0e\x32).protos.ModelUpdateRequest.TrainingStatus\x12\x11\n\tclient_id\x18\x05 \x01(\x05\x12\x12\n\nbatch_size\x18\x06 \x01(\x05\x12\x17\n\x0frequest_counter\x18\x07 \x01(\x05\x12\x1b\n\x13total_clients_count\x18\x08 \x01(\x05\x12\x1b\n\x13\x63lient_latest_clock\x18\t \x01(\x05\"+\n\x0eTrainingStatus\x12\r\n\tSUCCEEDED\x10\x00\x12\n\n\x06\x46\x41ILED\x10\x01\"c\n\x0cModelWeights\x12\x0e\n\x06layers\x18\x01 \x03(\x0c\x12\x13\n\x0bmodel_clock\x18\x02 \x01(\x05\x12\x11\n\tclient_id\x18\x03 \x01(\x05\x12\x1b\n\x13\x63lient_latest_clock\x18\t \x01(\x05\"\xc3\x01\n\x0cLayerWeights\x12:\n\x0bweight_type\x18\x01 \x01(\x0e\x32%.protos.LayerWeights.LayerWeightsType\x12\x15\n\rone_d_weights\x18\x02 \x03(\x01\x12\x15\n\rtwo_d_weights\x18\x03 \x03(\x0c\x12\x12\n\nlayer_rank\x18\x04 \x01(\x05\"5\n\x10LayerWeightsType\x12\x0b\n\x07UNKNOWN\x10\x00\x12\t\n\x05ONE_D\x10\x01\x12\t\n\x05TWO_D\x10\x02\"#\n\x10LayerTwoDWeights\x12\x0f\n\x07weights\x18\x01 \x03(\x01\x32\x91\x01\n\x10\x46\x65\x64\x65ratedService\x12J\n\x10TrainUpdateModel\x12\x1a.protos.ModelUpdateRequest\x1a\x1a.protos.ModelUpdateRequest\x12\x31\n\x11\x43hangeModeRestart\x12\r.protos.Empty\x1a\r.protos.EmptyB\x0bZ\t./;protosb\x06proto3'
  ,
  dependencies=[shared__pb2.DESCRIPTOR,])



_MODELUPDATEREQUEST_TRAININGSTATUS = _descriptor.EnumDescriptor(
  name='TrainingStatus',
  full_name='protos.ModelUpdateRequest.TrainingStatus',
  filename=None,
  file=DESCRIPTOR,
  create_key=_descriptor._internal_create_key,
  values=[
    _descriptor.EnumValueDescriptor(
      name='SUCCEEDED', index=0, number=0,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='FAILED', index=1, number=1,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=306,
  serialized_end=349,
)
_sym_db.RegisterEnumDescriptor(_MODELUPDATEREQUEST_TRAININGSTATUS)

_LAYERWEIGHTS_LAYERWEIGHTSTYPE = _descriptor.EnumDescriptor(
  name='LayerWeightsType',
  full_name='protos.LayerWeights.LayerWeightsType',
  filename=None,
  file=DESCRIPTOR,
  create_key=_descriptor._internal_create_key,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNKNOWN', index=0, number=0,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='ONE_D', index=1, number=1,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='TWO_D', index=2, number=2,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=595,
  serialized_end=648,
)
_sym_db.RegisterEnumDescriptor(_LAYERWEIGHTS_LAYERWEIGHTSTYPE)


_MODELUPDATEREQUEST = _descriptor.Descriptor(
  name='ModelUpdateRequest',
  full_name='protos.ModelUpdateRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='model_id', full_name='protos.ModelUpdateRequest.model_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='model_type', full_name='protos.ModelUpdateRequest.model_type', index=1,
      number=10, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='layers_weight', full_name='protos.ModelUpdateRequest.layers_weight', index=2,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='status', full_name='protos.ModelUpdateRequest.status', index=3,
      number=3, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='client_id', full_name='protos.ModelUpdateRequest.client_id', index=4,
      number=5, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='batch_size', full_name='protos.ModelUpdateRequest.batch_size', index=5,
      number=6, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='request_counter', full_name='protos.ModelUpdateRequest.request_counter', index=6,
      number=7, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='total_clients_count', full_name='protos.ModelUpdateRequest.total_clients_count', index=7,
      number=8, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='client_latest_clock', full_name='protos.ModelUpdateRequest.client_latest_clock', index=8,
      number=9, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _MODELUPDATEREQUEST_TRAININGSTATUS,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=42,
  serialized_end=349,
)


_MODELWEIGHTS = _descriptor.Descriptor(
  name='ModelWeights',
  full_name='protos.ModelWeights',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='layers', full_name='protos.ModelWeights.layers', index=0,
      number=1, type=12, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='model_clock', full_name='protos.ModelWeights.model_clock', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='client_id', full_name='protos.ModelWeights.client_id', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='client_latest_clock', full_name='protos.ModelWeights.client_latest_clock', index=3,
      number=9, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=351,
  serialized_end=450,
)


_LAYERWEIGHTS = _descriptor.Descriptor(
  name='LayerWeights',
  full_name='protos.LayerWeights',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='weight_type', full_name='protos.LayerWeights.weight_type', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='one_d_weights', full_name='protos.LayerWeights.one_d_weights', index=1,
      number=2, type=1, cpp_type=5, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='two_d_weights', full_name='protos.LayerWeights.two_d_weights', index=2,
      number=3, type=12, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='layer_rank', full_name='protos.LayerWeights.layer_rank', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _LAYERWEIGHTS_LAYERWEIGHTSTYPE,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=453,
  serialized_end=648,
)


_LAYERTWODWEIGHTS = _descriptor.Descriptor(
  name='LayerTwoDWeights',
  full_name='protos.LayerTwoDWeights',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='weights', full_name='protos.LayerTwoDWeights.weights', index=0,
      number=1, type=1, cpp_type=5, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=650,
  serialized_end=685,
)

_MODELUPDATEREQUEST.fields_by_name['status'].enum_type = _MODELUPDATEREQUEST_TRAININGSTATUS
_MODELUPDATEREQUEST_TRAININGSTATUS.containing_type = _MODELUPDATEREQUEST
_LAYERWEIGHTS.fields_by_name['weight_type'].enum_type = _LAYERWEIGHTS_LAYERWEIGHTSTYPE
_LAYERWEIGHTS_LAYERWEIGHTSTYPE.containing_type = _LAYERWEIGHTS
DESCRIPTOR.message_types_by_name['ModelUpdateRequest'] = _MODELUPDATEREQUEST
DESCRIPTOR.message_types_by_name['ModelWeights'] = _MODELWEIGHTS
DESCRIPTOR.message_types_by_name['LayerWeights'] = _LAYERWEIGHTS
DESCRIPTOR.message_types_by_name['LayerTwoDWeights'] = _LAYERTWODWEIGHTS
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ModelUpdateRequest = _reflection.GeneratedProtocolMessageType('ModelUpdateRequest', (_message.Message,), {
  'DESCRIPTOR' : _MODELUPDATEREQUEST,
  '__module__' : 'federated_pb2'
  # @@protoc_insertion_point(class_scope:protos.ModelUpdateRequest)
  })
_sym_db.RegisterMessage(ModelUpdateRequest)

ModelWeights = _reflection.GeneratedProtocolMessageType('ModelWeights', (_message.Message,), {
  'DESCRIPTOR' : _MODELWEIGHTS,
  '__module__' : 'federated_pb2'
  # @@protoc_insertion_point(class_scope:protos.ModelWeights)
  })
_sym_db.RegisterMessage(ModelWeights)

LayerWeights = _reflection.GeneratedProtocolMessageType('LayerWeights', (_message.Message,), {
  'DESCRIPTOR' : _LAYERWEIGHTS,
  '__module__' : 'federated_pb2'
  # @@protoc_insertion_point(class_scope:protos.LayerWeights)
  })
_sym_db.RegisterMessage(LayerWeights)

LayerTwoDWeights = _reflection.GeneratedProtocolMessageType('LayerTwoDWeights', (_message.Message,), {
  'DESCRIPTOR' : _LAYERTWODWEIGHTS,
  '__module__' : 'federated_pb2'
  # @@protoc_insertion_point(class_scope:protos.LayerTwoDWeights)
  })
_sym_db.RegisterMessage(LayerTwoDWeights)


DESCRIPTOR._options = None

_FEDERATEDSERVICE = _descriptor.ServiceDescriptor(
  name='FederatedService',
  full_name='protos.FederatedService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=688,
  serialized_end=833,
  methods=[
  _descriptor.MethodDescriptor(
    name='TrainUpdateModel',
    full_name='protos.FederatedService.TrainUpdateModel',
    index=0,
    containing_service=None,
    input_type=_MODELUPDATEREQUEST,
    output_type=_MODELUPDATEREQUEST,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='ChangeModeRestart',
    full_name='protos.FederatedService.ChangeModeRestart',
    index=1,
    containing_service=None,
    input_type=shared__pb2._EMPTY,
    output_type=shared__pb2._EMPTY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_FEDERATEDSERVICE)

DESCRIPTOR.services_by_name['FederatedService'] = _FEDERATEDSERVICE

# @@protoc_insertion_point(module_scope)
