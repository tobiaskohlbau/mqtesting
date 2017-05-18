import { Injectable } from '@angular/core';

@Injectable()
export class PostAuthService {
    addPostAuthListener(fn: EventListener): void {
        window.addEventListener('message', fn, false);
    }
}