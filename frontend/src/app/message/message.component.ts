import { Component, OnInit, OnDestroy } from '@angular/core';

import { Observable } from 'rxjs/Observable';
import "rxjs/add/observable/timer";
import { AnonymousSubscription } from "rxjs/Subscription";

import { Message } from './message.model';
import { MessageService } from './message.service';

import { AuthService } from '../auth/auth.service';

@Component({
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrls: ['./message.component.scss']
})
export class MessageComponent implements OnInit, OnDestroy {
  msgs: Message[];
  messageSubscription: AnonymousSubscription;
  timerSubscription: AnonymousSubscription;

  constructor(private messageService: MessageService, private authService: AuthService) {
  }

  ngOnInit() {
    this.getMessages();
  }

  ngOnDestroy() {
    if (this.messageSubscription) {
      this.messageSubscription.unsubscribe();
    }
    if (this.timerSubscription) {
      this.timerSubscription.unsubscribe();
    }
  }

  private getMessages() {
    this.messageSubscription = this.messageService.getMessages().subscribe(msgs => {
      if (this.msgs == null || this.msgs.length !== msgs.length) {
        this.msgs = msgs;
      }
      this.subscribeToMessages();
    });
  }

  private subscribeToMessages() {
    this.timerSubscription = Observable.timer(1000).first().subscribe(() => this.getMessages());
  }

  delete(id: string) {
    this.messageService.deleteMessage(id).subscribe(() => this.getMessages());
  }
}
