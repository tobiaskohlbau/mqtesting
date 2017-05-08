import { JsonInterface } from '../rest/json.interface';

export class Message implements JsonInterface {
  id: string;
  topic: string;
  payload: string;

  constructor(json?: Object) {
    if (json != null) {
      this.id = json['id'];
      this.topic = json['topic'];
      this.payload = json['payload']
    }
  }

  json(): string {
    return JSON.stringify(this);
  }
}
