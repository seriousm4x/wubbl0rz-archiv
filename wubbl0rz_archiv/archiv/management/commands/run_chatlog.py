import logging
import socket

from archiv.models import ChatMessage
from django.core.management.base import BaseCommand


class Command(BaseCommand):
    def handle(self, **options):
        c = Client()
        c.run()


class Client:
    def __init__(self) -> None:
        self.server = "irc.chat.twitch.tv"
        self.port = 6667
        self.nick = "justinfan00000"
        self.channel = "wubbl0rz"

        self.logger = logging.getLogger('wublog')
        self.logger.setLevel(logging.DEBUG)
        ch = logging.StreamHandler()
        ch.setLevel(logging.DEBUG)
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        ch.setFormatter(formatter)
        self.logger.addHandler(ch)

    def connect(self):
        self.logger.debug("connecting...")
        self.sock = socket.socket()
        self.sock.connect((self.server, self.port))
        self.send(f"NICK {self.nick}")
        self.send(f"JOIN #{self.channel}")
        self.send(f"CAP REQ :twitch.tv/tags twitch.tv/commands")

    def disconnect(self):
        self.sock.close()
        self.logger.debug("disconnected")

    def send(self, message):
        self.sock.send(f"{message}\n".encode("utf-8"))

    def receive(self):
        resp = self.sock.recv(2048).decode("utf-8").rstrip()
        if len(resp) == 0:
            return
        elif ":Your host is tmi.twitch.tv" in resp:
            self.logger.debug("connected")
            return
        elif resp.startswith("PING"):
            self.send("PONG")
            return
        elif resp.startswith(f":{self.nick}"):
            return
        self.logger.info(resp)

        ChatMessage.objects.create(raw=resp)

    def run(self):
        try:
            self.connect()
            while True:
                try:
                    self.receive()
                except Exception as e:
                    self.logger.error(e, exc_info=True)
                    self.disconnect()
                    self.connect()
        except KeyboardInterrupt:
            self.logger.debug("Quit by user")
            exit()
