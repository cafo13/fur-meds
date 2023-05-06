import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { AddVetAppointmentPage } from './add-vet-appointment.page';

describe('AddPetPage', () => {
  let component: AddVetAppointmentPage;
  let fixture: ComponentFixture<AddVetAppointmentPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddVetAppointmentPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(AddVetAppointmentPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
