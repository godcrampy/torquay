import random
import threading

from locust import HttpUser, TaskSet, task


class TokenTaskSet(TaskSet):
    received_tokens = set()
    token_lock = threading.Lock()
    hosts = ["http://localhost:8080", "http://localhost:8081"]

    @task
    def get_token(self):
        host = random.choice(self.hosts)
        with self.client.get(host + "/api/v1/token", catch_response=True) as response:
            if response.status_code == 200:
                token = response.text.strip()
                with self.token_lock:
                    if token in self.received_tokens:
                        response.failure(f"Duplicate token found: {token}")
                    else:
                        response.success()
                        self.received_tokens.add(token)
            else:
                response.failure(f"Failed to get token from {host}, status code: {response.status_code}")


class User(HttpUser):
    tasks = [TokenTaskSet]

    def on_start(self):
        self.host = None
