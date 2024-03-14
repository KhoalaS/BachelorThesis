from flask import Flask
from flask import request

app = Flask(__name__)


@app.route("/remove", methods=["POST"])
def hello_world():
    data = request.get_json(force=True)
    edge = data["edge"]
    current_rule = data["current_rule"]
    
    return "<p>Hello, World!</p>"
