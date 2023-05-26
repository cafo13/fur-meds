export enum AnimalSpecies {
  CAT = 'types.animal_species.cat',
  DOG = 'types.animal_species.dog',
  OTHER = 'types.animal_species.other',
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

export type PetFoodFrequency = {
  uuid: string;
  time: string;
};

export type PetFood = {
  uuid: string;
  name: string;
  dosage: string;
  frequencies: PetFoodFrequency[];
};

export type Pet = {
  uuid: string;
  name: string;
  species: AnimalSpecies | undefined;
  image: string;
  medicines?: PetMedicine[];
  foods?: PetFood[];
};
