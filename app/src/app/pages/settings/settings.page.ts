import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { IonicModule } from '@ionic/angular';
import { TranslocoRootModule } from 'src/app/transloco-root.module';
import { environment } from 'src/environments/environment';
import { TranslocoService } from '@ngneat/transloco';

type Language = Record<'code' | 'name' | 'shorthand', string>;
@Component({
  selector: 'app-settings',
  templateUrl: 'settings.page.html',
  styleUrls: ['settings.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoRootModule],
})
export class SettingsPage {
  useDarkTheme: boolean;
  selectedLanguageCode: string;
  languagesList: Array<Language> = [
    {
      code: 'en',
      name: 'English',
      shorthand: 'ENG',
    },
    {
      code: 'de',
      name: 'German',
      shorthand: 'GER',
    },
  ];

  constructor(protected transloco: TranslocoService) {
    const appSettingsUseDarkTheme = localStorage.getItem('useDarkTheme');
    const deviceSettingsUseDarkTheme = window.matchMedia(
      '(prefers-color-scheme: dark)'
    ).matches;
    if (appSettingsUseDarkTheme !== null) {
      this.useDarkTheme = JSON.parse(appSettingsUseDarkTheme);
    } else {
      this.useDarkTheme = deviceSettingsUseDarkTheme;
    }
    document.body.classList.toggle('dark', this.useDarkTheme);

    this.selectedLanguageCode = this.transloco.getActiveLang();
  }

  toggleDarkTheme(event: any): void {
    const useDarkTheme = event.detail.checked;
    document.body.classList.toggle('dark', useDarkTheme);
    this.useDarkTheme = useDarkTheme;
    localStorage.setItem('useDarkTheme', useDarkTheme);
  }

  handleLanguageChange(event: any): void {
    const languageCode = event.detail.value;
    this.transloco.setActiveLang(languageCode);
    this.selectedLanguageCode = languageCode;
    localStorage.setItem('language', languageCode);
  }

  getAppVersion() {
    return environment.version;
  }
}
