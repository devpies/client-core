function IsNew(user, context, callback) {
  const namespace = "https://client-domain.com/claims/";
  context.idToken[namespace + "is_new"] = context.stats.loginsCount === 1;
  context.accessToken[namespace + "is_new"] = context.stats.loginsCount === 1;
  callback(null, user, context);
}
