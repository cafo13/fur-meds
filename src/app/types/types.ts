export type PetMedicine = {
  name: string;
};

export type PetFood = {
  name: string;
};

export type Pet = {
  name: string;
  image: string;
  medicines: PetMedicine[];
  foods: PetFood[];
};
