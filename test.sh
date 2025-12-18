#!/bin/bash

# Arrête le script dès qu'une commande échoue
set -e

# Fonction pour afficher les erreurs
error_exit() {
    echo ""
    echo "❌ ERREUR à l'étape: $1"
    echo "Le script s'est arrêté. Corrige l'erreur avant de continuer."
    exit 1
}

echo "=== Étape 2: Vérification des tests initiaux ==="
go test ./... || error_exit "Étape 2 - Les tests initiaux échouent"

echo ""
echo "=== Étape 3: Mise à jour de la version Go ==="
GO_VERSION="1.25.5"
echo "Mise à jour vers Go $GO_VERSION"
go mod edit -go=$GO_VERSION || error_exit "Étape 3 - Impossible de mettre à jour la version Go"
echo "⚠️  N'oublie pas de vérifier manuellement les fichiers .github/workflows/*.yml si nécessaire"

echo ""
echo "=== Étape 4: Mise à jour des dépendances Go ==="
go get -u -t ./... || error_exit "Étape 4 - Échec de la mise à jour des dépendances"
go mod tidy || error_exit "Étape 4 - Échec de go mod tidy"

if [ -d "vendor" ]; then
    echo "Mise à jour du répertoire vendor..."
    go mod vendor || error_exit "Étape 4 - Échec de go mod vendor"
else
    echo "Pas de répertoire vendor détecté, étape omise"
fi

echo ""
echo "=== Étape 5: Lancement du linter ==="
if [ -f ".golangci.yml" ]; then
    echo "Fichier .golangci.yml trouvé"
    if ! golangci-lint run --fix --config .golangci.yml; then
        echo "Tentative de migration de la config..."
        golangci-lint migrate --skip-validation || error_exit "Étape 5 - Échec de la migration golangci-lint"
        golangci-lint run --fix --config .golangci.yml || error_exit "Étape 5 - Le linter a trouvé des erreurs à corriger"
    fi
else
    echo "Pas de fichier .golangci.yml, utilisation de la config par défaut"
    golangci-lint run --fix || error_exit "Étape 5 - Le linter a trouvé des erreurs à corriger"
fi

echo ""
echo "=== Étape 6: Relance des tests ==="
go test ./... || error_exit "Étape 6 - Les tests échouent après les modifications"

echo ""
echo "=== Étape 7: Build ==="
if [ -f ".goreleaser.yml" ] || [ -f ".goreleaser.yaml" ]; then
    echo "GoReleaser détecté, lancement du build..."
    goreleaser release --clean --skip=publish --skip=validate --skip=sign || error_exit "Étape 7 - Échec du build GoReleaser"
else
    echo "Pas de GoReleaser, build Go classique..."
    if [ -d "vendor" ]; then
        go build -mod=vendor || error_exit "Étape 7 - Échec du build avec vendor"
    else
        go build || error_exit "Étape 7 - Échec du build"
    fi
fi

echo ""
echo "✅ Toutes les étapes sont terminées avec succès !"
echo "Tu peux maintenant commiter tes changements (étape 8)"
