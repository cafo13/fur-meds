import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { AddMedicinePage } from './add-medicine.page';

describe('AddPetPage', () => {
  let component: AddMedicinePage;
  let fixture: ComponentFixture<AddMedicinePage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddMedicinePage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(AddMedicinePage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
