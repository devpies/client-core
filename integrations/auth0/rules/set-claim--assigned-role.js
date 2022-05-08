function SetAssignedRoles(user, context, callback) {
    const namespace = 'https://client-domain.com/claims/roles';
    const assignedRoles = (context.authorization || {}).roles;

    user.user_metadata = user.user_metadata || {};

    if(!assignedRoles.includes(user.user_metadata.role)) {
        assignedRoles.push(user.user_metadata.role)
    }

    context.idToken[namespace] = assignedRoles;
    context.accessToken[namespace] = assignedRoles;

    callback(null, user, context);
}
