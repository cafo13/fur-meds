export enum PetSpecies {
  CAT = 'Cat',
  DOG = 'Dog',
}

export type PetMedicine = {
  name: string;
};

export type PetFood = {
  name: string;
};

export type VetAppointments = {
  name: string;
};

export type Pet = {
  name: string;
  species?: PetSpecies;
  image?: string;
  medicines?: PetMedicine[];
  foods?: PetFood[];
  vetAppointments?: VetAppointments[];
};
