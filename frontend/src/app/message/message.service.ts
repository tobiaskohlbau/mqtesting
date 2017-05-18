import { Injectable } from '@angular/core';

import { Observable } from "rxjs/Observable";
import "rxjs/add/operator/map";
import { Message } from './message.model';
import { RestService } from '../rest/rest.service';

@Injectable()
export class MessageService {
  constructor(private rest: RestService) {
    rest.setURL('https://localhost/api/messages');
  }

  getMessage(id: string): Observable<Message> {
        return this.rest.read(id).map(msg => new Message(msg));
  }

  getMessages(): Observable<Array<Message>> {
    return this.rest.read().map(msgs => msgs.map(msg => new Message(msg)));
  }

  deleteMessage(id: string): Observable<Message> {
    return this.rest.delete(id).map(msg => new Message(msg));
  }
}
