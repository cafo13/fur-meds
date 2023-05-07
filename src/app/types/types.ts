export enum PetSpecies {
  CAT = 'Cat',
  DOG = 'Dog',
}

export type PetMedicineFrequency = {
  uuid: string;
  time: string;
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

export type Pet = {
  uuid: string;
  name: string;
  species?: PetSpecies;
  image?: string;
  medicines?: PetMedicine[];
  foods?: PetFood[];
};
