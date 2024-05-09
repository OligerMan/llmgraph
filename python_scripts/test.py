import pika
import json

connection = pika.BlockingConnection()
channel = connection.channel()

method_frame, properties, body = next(channel.consume('start_node'))
channel.basic_ack(method_frame.delivery_tag)

data = json.loads(body.decode("utf-8"))
print("Request: ", data)
print("Question: ", data[0]["value_info"])

answer = [{'wrapper_id': 'string', 'value_type': 'string', 'value_info': 'Four'}]
answer = bytes(json.dumps(answer), "utf-8")
print("Sending data: ", answer)
channel.basic_publish(exchange='', routing_key='start_node_response', body=answer)