export class AuthenticationError extends Error {
  constructor(message = "Unable to authenticate") {
    super(message);
  }
}
