import { Component, OnInit } from '@angular/core';

import { Message } from './message.model';
import { MessageService } from './message.service';

@Component({
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrls: ['./message.component.scss']
})
export class MessageComponent implements OnInit {
  msgs: Message[];

  constructor(private messageService: MessageService) { }

  ngOnInit() {
    this.pollMessages();
  }

  pollMessages() {
    this.messageService.pollMessages().subscribe(msgs => this.msgs = msgs);
  }

  getMessages() {
    this.messageService.getMessages().subscribe(msgs => this.msgs = msgs);
  }

  delete(id: string) {
    this.messageService.deleteMessage(id).subscribe(() => this.getMessages());
  }
}
