# Include gomk if it's been checked-out: git submodule update --init
-include gomk/main.mk
-include local/Makefile

PKGS := $(filter-out ./testhelper,$(PKGS))

superclean: superclean-default
ifeq ($(unameS),windows)
	@remove-item -force -recurse ./git/testdata
	@remove-item -force -recurse ./testdata
else
	@rm -f -r ./git/testdata ./testdata
endif
