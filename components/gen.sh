#!/bin/bash
# Minimum required version
required_version="4.0"

# Get the current version
current_version="${BASH_VERSION%%.*}"

# Compare versions
if (( $(echo "$current_version < $required_version" | bc -l) )); then
    echo "Bash version $current_version is less than required version $required_version."
    exit 1
fi

(
echo '// this file is autogenerated by gen.sh DO NOT EDIT.'
echo 'package components'
echo ''
echo '//go:generate ./gen.sh'
echo ''

echo 'import ('
echo '	"context"'
echo ''
echo '	"github.com/gosthome/gosthome/core/component"'
echo '	"github.com/gosthome/gosthome/core/registry"'
ls -l | awk '/^d/ {print $9}' | sed 's|^|\t"github.com/gosthome/gosthome/components/|g;s|$|"|g'
echo ')'
echo ''

ls -l | awk '/^d/ {print $9}' | while read dn; do
echo ''
echo ''
echo "type ${dn}Component struct {"
echo "}"
echo ''
echo "func (${dn}Component) Config() *component.ConfigDecoder {"
echo "	return component.NewConfigDecoder(${dn}.NewConfig())"
echo "}"
echo "func (${dn}Component) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {"
echo "	${dn}Cfg := cfg.(*${dn}.Config)"
echo "	return ${dn}.New(ctx, ${dn}Cfg)"
echo "}"
find -mindepth 2 -maxdepth 2 -type f -name domain.go -printf '%h\n' | sed 's|\./||g' | while read pn; do
if [ -f "${dn}/${dn}${pn}.go" ]; then
echo "type ${dn}${pn^}EntityComponent struct {"
echo "}"
echo ''
echo "func (${dn}${pn^}EntityComponent) Config() *component.ConfigDecoder {"
echo "	return component.NewConfigDecoder(${dn}.New${pn^}Config())"
echo "}"
echo "func (${dn}${pn^}EntityComponent) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {"
echo "	${dn}${pn^}Cfg := cfg.(*${dn}.${pn^}Config)"
echo "	return ${dn}.New${pn^}(ctx, ${dn}${pn^}Cfg)"
echo "}"
echo "func (${dn}Component) ${pn^}Platform() component.Declaration {"
echo "	return &${dn}${pn^}EntityComponent{}"
echo "}"
fi
done
done

echo "var ("
ls -l | awk '/^d/ {print $9}' | while read dn; do
if [ -f "${dn}/name.go" ]; then
echo "	COMPONENT_KEY_${dn^^} = ${dn}.COMPONENT_KEY"
else
echo "	COMPONENT_KEY_${dn^^} = \"${dn}\""
fi
done
echo ")"

# echo "func RegisterAll(cr *registry.Registry) (err error) {"
# ls -l | awk '/^d/ {print $9}' | while read dn; do
# echo "	if err = cr.Register(COMPONENT_KEY_${dn^^}, ${dn}Component{}); err != nil {"
# echo "		return err"
# echo "	}"
# done
# echo "	return nil"
# echo "}"

echo "var("
ls -l | awk '/^d/ {print $9}' | while read dn; do
echo "	_ = registry.RegisterDefaultComponent(COMPONENT_KEY_${dn^^}, ${dn}Component{})"

done
echo ")"

)> $GOFILE

gofumpt -w -extra $GOFILE
