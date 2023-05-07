import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { SetMedicinePage } from './set-medicine.page';

describe('SetMedicinePage', () => {
  let component: SetMedicinePage;
  let fixture: ComponentFixture<SetMedicinePage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SetMedicinePage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(SetMedicinePage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
