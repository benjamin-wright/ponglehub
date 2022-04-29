export interface UserData {
  name: string
}

export class PongleStorage {
  private storage: Storage;

  constructor(storage: Storage) {
    this.storage = storage;
  }

  isLoggedIn(): boolean {
    return !!this.storage.getItem('userData');
  }

  setUserData(data: UserData): void {
    this.storage.setItem('userData', JSON.stringify(data));
  }

  getUserData(): UserData {
    const userString = this.storage.getItem('userData');
    if (userString == null) {
      throw new Error('userData not found');
    }

    const parsed = JSON.parse(userString);
    if (typeof parsed.name !== "string") {
      throw new Error(`parsing error, name property not found on "${userString}"`);
    }
  
    return {
      name: parsed.name
    };
  }

  clear(): void {
    this.storage.clear();
  }
}
