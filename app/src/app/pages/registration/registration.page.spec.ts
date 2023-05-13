import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { RegistrationPage } from './registration.page';

describe('RegistrationPage', () => {
  let component: RegistrationPage;
  let fixture: ComponentFixture<RegistrationPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RegistrationPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(RegistrationPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
