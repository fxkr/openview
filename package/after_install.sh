set -e

OPENVIEW_USER="openview"
OPENVIEW_GROUP="openview"

# Create group
if ! getent group "${OPENVIEW_GROUP}" > /dev/null 2>&1 ; then
    addgroup \
        --system \
        --quiet \
        "${OPENVIEW_GROUP}"
fi

# Create user
if ! id "${OPENVIEW_USER}" > /dev/null 2>&1 ; then
    adduser \
        --system \
        --home /var/lib/openview \
        --no-create-home \
        --ingroup "${OPENVIEW_GROUP}" \
        --disabled-password \
        --shell /bin/false \
        "${OPENVIEW_USER}"
fi

# Apply permissions
chown "${OPENVIEW_USER}":"${OPENVIEW_USER}" "/var/cache/openview"