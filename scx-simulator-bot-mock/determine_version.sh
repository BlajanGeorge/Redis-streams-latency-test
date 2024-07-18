if [[ $(git rev-parse --abbrev-ref HEAD) == "main" ]]; then
  echo "$(mvn help:evaluate -Dexpression=project.version -q -DforceStdout).$(git rev-list --count HEAD)"
else
  echo "$(mvn help:evaluate -Dexpression=project.version -q -DforceStdout).$(git rev-list --count HEAD)-$(git rev-parse --abbrev-ref HEAD)-$(git rev-parse --short HEAD)" | sed -e "s|/|-|"
fi
