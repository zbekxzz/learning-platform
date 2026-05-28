import { Injectable } from '@angular/core';
import { MessageService } from 'primeng/api';

@Injectable({
  providedIn: 'root'
})
export class CustomMessageService {
  constructor(private messageService: MessageService) {}

  private showTranslatedMessage(severity: string, summaryKey: string, detailKey: string) {
    this.messageService.add({
      severity,
      summary: summaryKey,
      detail: detailKey,
      life: 3000
    });
  }

  showErrorMessage(errorMessage: string): void {
    this.showTranslatedMessage('error', 'Қате', errorMessage);
  }

  showInfoMessage(errorMessage: string): void {
    this.showTranslatedMessage('info', 'Ескерту', errorMessage);
  }

  showSuccessMessage(errorMessage: string): void {
    this.showTranslatedMessage('success', 'Сәтті', errorMessage);
  }
}
