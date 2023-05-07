export enum PetSpecies {
  CAT = 'Cat',
  DOG = 'Dog',
}

export type PetMedicineFrequencyTime = {
  hour: number;
  minute: number;
};

export type PetMedicineFrequency = {
  time: PetMedicineFrequencyTime;
  everyDays: number;
};

export type PetMedicine = {
  uuid: string;
  name: string;
  dosage: string;
  frequencies: PetMedicineFrequency[];
};

export type PetFood = {
  uuid: string;
  name: string;
};

export type PetVetAppointment = {
  uuid: string;
  name: string;
};

export type Pet = {
  uuid: string;
  name: string;
  species?: PetSpecies;
  image?: string;
  medicines?: PetMedicine[];
  foods?: PetFood[];
  vetAppointments?: PetVetAppointment[];
};
